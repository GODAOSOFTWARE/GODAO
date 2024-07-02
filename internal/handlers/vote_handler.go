package handlers

import (
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	"dao_vote/internal/models"
	"dao_vote/internal/services"
	"dao_vote/internal/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// CreateVoteHandler обрабатывает POST /votes запрос для создания нового голосования
func CreateVoteHandler(c *gin.Context) {
	logrus.Info("CreateVoteHandler started")

	// Получение токена из заголовка
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.JSONResponse(c, http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		logrus.Error("Authorization token is missing")
		return
	}

	// Запрос информации о текущем пользователе
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://backend.ddapps.io/api/v1/auth/me?with_user_information=1", nil)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to create request to auth/me: %v", err)
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to get response from auth/me: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to read response body from auth/me: %v", err)
		return
	}

	var userMeResp UserMeResponse
	if err := json.Unmarshal(body, &userMeResp); err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to unmarshal response from auth/me: %v", err)
		return
	}
	voter := userMeResp.Data.Wallet
	logrus.Infof("Retrieved voter wallet: %s", voter)

	// Получение данных из формы
	var vote models.VoteWithoutID
	vote.Title = c.PostForm("title")
	vote.Subtitle = c.PostForm("subtitle")
	vote.Description = c.PostForm("description")
	vote.Voter = voter
	vote.Choice = c.PostForm("choice")
	logrus.Infof("Form data received: %+v", vote)

	// Валидация данных голосования
	if err := validate.Struct(vote); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("Validation error: %v", err)
		return
	}
	logrus.Info("Vote data validated")

	// Создание мнемонической фразы для кошелька
	mnemonic, err := wallet.NewMnemonic(256, "")
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Не удалось создать мнемоническую фразу"})
		logrus.Errorf("Failed to create mnemonic: %v", err)
		return
	}

	// Создание аккаунта из мнемонической фразы
	account, err := wallet.NewAccountFromMnemonic(mnemonic)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Не удалось создать аккаунт"})
		logrus.Errorf("Failed to create account from mnemonic: %v", err)
		return
	}

	// Логирование мнемонической фразы и адреса кошелька
	logrus.Infof("Mnemonic: %s", mnemonic.Words())
	logrus.Infof("Wallet Address: %s", account.Address())

	voteWithID := models.Vote{
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
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to get vote strength: %v", err)
		return
	}
	voteWithID.VotePower = votePower
	logrus.Infof("Vote power obtained: %f", float64(votePower)) // Приведение к float64 для корректного форматирования

	// Сохранение голосования в базе данных
	id, err := services.CreateVote(voteWithID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		logrus.Errorf("Failed to create vote: %v", err)
		return
	}
	voteWithID.ID = id
	logrus.Infof("Vote saved with ID: %d", id)

	utils.JSONResponse(c, http.StatusCreated, voteWithID)
	logrus.Info("CreateVoteHandler completed successfully")
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
		logrus.Errorf("Vote not found: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, vote)
	logrus.Infof("Vote retrieved successfully: %+v", vote)
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
	logrus.Infof("Vote deleted successfully: %d", id)
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
		logrus.Errorf("Vote not found: %v", err)
		return
	}
	logrus.Infof("Vote retrieved successfully: %+v", vote)

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
	logrus.Infof("Vote power for voter %s: %d", vote.Voter, votePower)

	// Получение данных голоса пользователя
	var userVote models.UserVote
	if err := c.BindJSON(&userVote); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		logrus.Errorf("Invalid request body: %v", err)
		return
	}
	logrus.Infof("User vote data received: %+v", userVote)

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

	// Ответ клиенту
	utils.JSONResponse(c, http.StatusCreated, userVote)
	logrus.Info("AddUserVoteHandler completed successfully")
}

// GetUserVotesHandler обрабатывает GET /votes/:id/votes запрос для получения всех голосов пользователей для голосования
func GetUserVotesHandler(c *gin.Context) {
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid VoteID"})
		logrus.Errorf("Invalid VoteID: %v", err)
		return
	}

	userVotes, err := services.GetUserVotes(voteID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to get user votes"})
		logrus.Errorf("Failed to get user votes: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, userVotes)
	logrus.Infof("User votes retrieved successfully: %+v", userVotes)
}
