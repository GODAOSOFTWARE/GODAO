// Package handlers Обработчик голосований пользователей

package handlers

import (
	"bitbucket.org/decimalteam/dsc-go-sdk/wallet"
	"bytes"
	"dao_vote/back-end/models"
	"dao_vote/back-end/services"
	"dao_vote/back-end/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// GetVotingResultsByWallet обрабатывает GET запрос результатов голосования по кошельку
func GetVotingResultsByWallet(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		// Получаем параметр wallet_address из запроса
		walletAddress := c.Query("wallet_address")
		if walletAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_address is required"})
			return nil
		}

		// Получаем параметр offset из запроса
		offsetStr := c.Query("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0 // Устанавливаем значение по умолчанию, если параметр отсутствует или некорректен
		}

		// Вызов сервиса для получения результатов голосования команды DAO
		apiResponse, err := services.FetchVoteResults(walletAddress, offset)
		if err != nil {
			// Возвращает ошибку, если не удалось получить результаты голосования
			return err
		}

		// Подготовка результатов голосования для отправки в ответе
		voteResults := services.PrepareVoteResults(apiResponse)
		// Возвращает успешный ответ с результатами голосования
		c.JSON(http.StatusOK, voteResults)
		return nil
	})
}

// CreateVoteHandler обрабатывает запрос на создание голооования
func CreateVoteHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		logrus.Info("CreateVoteHandler started")

		// Получение токена из заголовка
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			logrus.Error("Authorization token is missing")
			return nil
		}

		// Запрос информации о текущем пользователе
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://backend.ddapps.io/api/v1/auth/me?with_user_information=1", nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var userMeResp UserMeResponse
		if err := json.Unmarshal(body, &userMeResp); err != nil {
			return err
		}
		voter := userMeResp.Data.Wallet
		logrus.Infof("Retrieved voter wallet: %s", voter)

		// Получение данных из формы
		var vote models.NewVote
		vote.Title = c.PostForm("title")
		vote.Subtitle = c.PostForm("subtitle")
		vote.Description = c.PostForm("description")
		vote.Voter = voter
		vote.Choice = c.PostForm("choice")
		logrus.Infof("Form data received: %+v", vote)

		// Валидация данных голосования
		if err := validate.Struct(vote); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			logrus.Errorf("Validation error: %v", err)
			return nil
		}
		logrus.Info("VoteInfo data validated")

		// Генерация мнемонической фразы для кошелька
		mnemonicObject, err := wallet.NewMnemonic("")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать мнемоническую фразу"})
			logrus.Errorf("Failed to create mnemonic: %v", err)
			return nil
		}
		mnemonic := mnemonicObject.Words()
		fmt.Println("Generated Mnemonic Words:", mnemonic)

		// Создание аккаунта из сгенерированной мнемонической фразы
		account, err := wallet.NewAccountFromMnemonicWords(mnemonic, "")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать аккаунт"})
			logrus.Errorf("Failed to create account from mnemonic: %v", err)
			return nil
		}
		fmt.Println("Generated Address from New Mnemonic:", account.Address())

		// Логирование мнемонической фразы и адреса кошелька
		logrus.Infof("Mnemonic: %s", mnemonic)
		logrus.Infof("Wallet Address: %s", account.Address())

		voteWithID := models.VoteInfo{
			Title:         vote.Title,
			Subtitle:      vote.Subtitle,
			Description:   vote.Description,
			Voter:         vote.Voter,
			Choice:        vote.Choice,
			WalletAddress: account.Address(),
		}

		// Получение силы голоса для голосующего
		votePower, err := services.GetVoteStrength(vote.Voter)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			logrus.Errorf("Failed to get vote strength: %v", err)
			return nil
		}
		voteWithID.VotePower = votePower
		logrus.Infof("VoteInfo power obtained: %f", float64(votePower)) // Приведение к float64 для корректного форматирования

		// Сохранение голосования в базе данных
		id, err := services.CreateVote(voteWithID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			logrus.Errorf("Failed to create vote: %v", err)
			return nil
		}
		voteWithID.ID = id
		logrus.Infof("VoteInfo saved with ID: %d", id)

		c.JSON(http.StatusCreated, voteWithID)
		logrus.Info("CreateVoteHandler completed successfully")
		return nil
	})
}

// GetVoteHandler получает голосование по ID
func GetVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoteID"})
		logrus.Errorf("Invalid VoteID: %v", err)
		return
	}

	vote, err := services.GetVote(id)
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		logrus.Errorf("VoteInfo not found: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, vote)
	logrus.Infof("VoteInfo retrieved successfully: %+v", vote)
}

// DeleteVoteHandler обрабатывает DELETE /votes/:id запрос для удаления голосования
func DeleteVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoteID"})
		logrus.Errorf("Invalid VoteID: %v", err)
		return
	}

	if err := services.DeleteVote(id); err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to delete vote: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusNoContent, gin.H{})
	logrus.Infof("VoteInfo deleted successfully: %d", id)
}

// AddUserVoteHandler добавляет голос пользователя к голосованию
func AddUserVoteHandler(c *gin.Context) {
	logrus.Info("AddUserVoteHandler started")

	// Проверка ID голосования
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid VoteID"})
		logrus.Errorf("Invalid VoteID: %v", err)
		return
	}
	logrus.Infof("VoteID: %d", voteID)

	// Получение данных голосования
	vote, err := services.GetVote(voteID)
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		logrus.Errorf("VoteInfo not found: %v", err)
		return
	}
	logrus.Infof("VoteInfo retrieved successfully: %+v", vote)

	// Сохранение кошелька голосования в памяти
	walletAddress := vote.WalletAddress
	logrus.Infof("Wallet address for the vote: %s", walletAddress)

	// Получение силы голоса создателя голосования
	votePower, err := services.GetVoteStrength(vote.Voter)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Error determining vote strength"})
		logrus.Errorf("Error determining vote strength: %v", err)
		return
	}
	logrus.Infof("VoteInfo power for voter %s: %d", vote.Voter, votePower)

	// Получение данных голоса пользователя
	var userVote models.UserVote
	if err := c.BindJSON(&userVote); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		logrus.Errorf("Invalid request body: %v", err)
		return
	}

	userVote.VoteID = voteID
	userVote.VotePower = votePower
	userVote.Voter = vote.Voter // Заполняем поле Voter в структуре UserVote

	// Сохранение голоса пользователя в базе данных
	id, err := services.AddUserVote(userVote)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to add vote"})
		logrus.Errorf("Failed to add vote: %v", err)
		return
	}

	userVote.VoterID = id
	logrus.Infof("User vote added successfully: %+v", userVote)

	// Инициация вывода средств
	logrus.Info("Initiating withdrawal")

	amount := 1 // Установка количества средств
	logrus.Infof("Amount for withdrawal: %d", amount)

	token := c.GetHeader("Authorization")
	if token == "" {
		logrus.Error("Authorization token is missing")
		utils.JSONResponse(c, http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}
	logrus.Infof("Authorization token: %s", token)

	withdrawReq := models.WithdrawRequest{
		Amount:  float64(amount),
		Address: walletAddress,
	}
	logrus.Infof("Withdraw request data: %+v", withdrawReq)

	// Вызов функции вывода средств
	response, err := initiateWithdrawal(withdrawReq, token)
	if err != nil {
		logrus.Errorf("Failed to initiate withdrawal: %v", err)
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to initiate withdrawal"})
		return
	}
	logrus.Infof("Withdraw response: %+v", response)

	transactionID := response.Data.TransactionID
	if transactionID == 0 {
		logrus.Error("Invalid transaction ID in response")
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Invalid transaction ID in response"})
		return
	}

	// Запрос хэша транзакции
	logrus.Infof("Waiting before requesting transaction hash")
	time.Sleep(5 * time.Second) // Задержка 5 секунд

	hashResponse, err := getTransactionHash(transactionID, token)
	if err != nil || hashResponse.Data.Hash == "" {
		logrus.Errorf("Failed to retrieve transaction hash: %v", err)
		logrus.Infof("Retrying to retrieve transaction hash after 5 seconds delay")
		time.Sleep(5 * time.Second) // Повторная задержка 5 секунд

		hashResponse, err = getTransactionHash(transactionID, token)
		if err != nil {
			logrus.Errorf("Failed to retrieve transaction hash on second attempt: %v", err)
			utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return
		}
	}
	logrus.Infof("Transaction hash response: %+v", hashResponse)

	transactionHash := hashResponse.Data.Hash
	logrus.Infof("Transaction hash: %s", transactionHash)

	if transactionHash == "" {
		logrus.Error("Invalid transaction hash in response")
		utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	// Ответ клиенту
	utils.JSONResponse(c, http.StatusOK, gin.H{
		"message":          "Withdrawal successful",
		"transaction_id":   transactionID,
		"transaction_hash": transactionHash,
	})
	logrus.Info("AddUserVoteHandler completed successfully")
}

// initiateWithdrawal вызывает существующую функцию вывода средств
func initiateWithdrawal(req models.WithdrawRequest, token string) (WithdrawResponse, error) {
	client := &http.Client{}
	jsonData, err := json.Marshal(req)
	if err != nil {
		logrus.Errorf("Failed to marshal JSON: %v", err)
		return WithdrawResponse{}, err
	}

	logrus.Infof("Initiating HTTP request for withdrawal")
	httpReq, err := http.NewRequest("POST", "https://backend.ddapps.io/api/v1/withdraw", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("Failed to create new HTTP request: %v", err)
		return WithdrawResponse{}, err
	}

	httpReq.Header.Set("Authorization", token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Логгирование содержимого HTTP-запроса
	logrus.Infof("HTTP Request - URL: %s", httpReq.URL)
	logrus.Infof("HTTP Request - Method: %s", httpReq.Method)
	logrus.Infof("HTTP Request - Headers: %v", httpReq.Header)
	logrus.Infof("HTTP Request - Body: %s", jsonData)

	logrus.Infof("Sending request to external API for withdrawal")
	resp, err := client.Do(httpReq)
	if err != nil {
		logrus.Errorf("Failed to send request to external API: %v", err)
		return WithdrawResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body: %v", err)
		return WithdrawResponse{}, err
	}

	logrus.Infof("Response from external API: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		return WithdrawResponse{}, fmt.Errorf("error from external API: %s", string(body))
	}

	var response WithdrawResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logrus.Errorf("Failed to unmarshal response JSON: %v", err)
		return WithdrawResponse{}, err
	}

	return response, nil
}

// getTransactionHash выполняет запрос для получения хэша транзакции
func getTransactionHash(transactionID int, token string) (TransactionHashResponse, error) {
	client := &http.Client{}
	hashReqURL := fmt.Sprintf("https://backend.ddapps.io/api/v1/transactions/%d", transactionID)
	hashReq, err := http.NewRequest("GET", hashReqURL, nil)
	if err != nil {
		logrus.Errorf("Failed to create new request for transaction hash: %v", err)
		return TransactionHashResponse{}, err
	}
	hashReq.Header.Set("Authorization", token)
	hashReq.Header.Set("Accept", "application/json")

	logrus.Infof("Sending request to external API for transaction hash: %v", hashReqURL)
	logrus.Infof("Request URL: %v", hashReqURL)
	logrus.Infof("Request method: %v", hashReq.Method)
	logrus.Infof("Request headers: %v", hashReq.Header)
	logrus.Infof("Request ID: %d", transactionID)

	time.Sleep(5 * time.Second) // Задержка 5 секунд
	hashResp, err := client.Do(hashReq)
	if err != nil {
		logrus.Errorf("Failed to send request to external API for transaction hash: %v", err)
		return TransactionHashResponse{}, err
	}
	defer hashResp.Body.Close()

	hashBody, err := ioutil.ReadAll(hashResp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body for transaction hash: %v", err)
		return TransactionHashResponse{}, err
	}

	logrus.Infof("Response from external API for transaction hash: %s", string(hashBody))
	if hashResp.StatusCode != http.StatusOK {
		return TransactionHashResponse{}, fmt.Errorf("error from external API for transaction hash: %s", string(hashBody))
	}

	var hashResponse TransactionHashResponse
	if err := json.Unmarshal(hashBody, &hashResponse); err != nil {
		logrus.Errorf("Failed to unmarshal response JSON for transaction hash: %v", err)
		return TransactionHashResponse{}, err
	}

	return hashResponse, nil
}

// Структуры для обработки ответов от внешнего API
type WithdrawResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Data    struct {
		TransactionID int `json:"transaction_id"`
	} `json:"data"`
}

type TransactionHashResponse struct {
	Data struct {
		ID             int    `json:"id"`
		UserID         int    `json:"user_id"`
		WalletID       int    `json:"wallet_id"`
		Amount         string `json:"amount"`
		Payload        string `json:"payload"`
		DecimalPayload string `json:"decimal_payload"`
		Coin           string `json:"coin"`
		Hash           string `json:"hash"`
		Result         string `json:"result"`
		Complete       bool   `json:"complete"`
		Success        bool   `json:"success"`
		CreatedAt      string `json:"created_at"`
		UpdatedAt      string `json:"updated_at"`
		Type           int    `json:"type"`
		Calculated     int    `json:"calculated"`
	} `json:"data"`
}

// GetUserVotesHandler обрабатывает GET /votes/:id/votes запрос для получения всех голосов пользователей для голосования
func GetUserVotesHandler(c *gin.Context) {
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid VoteID"})
		logrus.Errorf("Invalid VoteID: %v", err)
		return
	}

	voteResults, err := services.FetchVotes(voteID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to get vote results"})
		logrus.Errorf("Failed to get vote results: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, voteResults)
	logrus.Infof("VoteInfo results retrieved successfully: %+v", voteResults)
}
