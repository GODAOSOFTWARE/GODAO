package handlers

import (
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	"dao_vote/internal/models"
	"dao_vote/internal/services"
	"dao_vote/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// CreateVoteHandler обрабатывает POST /votes запрос для создания нового голосования
func CreateVoteHandler(c *gin.Context) {
	var vote models.VoteWithoutID
	vote.Title = c.PostForm("title")
	vote.Subtitle = c.PostForm("subtitle")
	vote.Description = c.PostForm("description")
	vote.Voter = c.PostForm("voter")
	vote.Choice = c.PostForm("choice")

	// Проверяем правильность заполнения полей
	if err := validate.Struct(vote); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация новой мнемонической фразы
	mnemonic, err := wallet.NewMnemonic(256, "")
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Не удалось создать мнемоническую фразу"})
		return
	}

	// Создание нового аккаунта из мнемонической фразы
	account, err := wallet.NewAccountFromMnemonic(mnemonic)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Не удалось создать аккаунт"})
		return
	}

	// Логирование мнемонической фразы и адреса кошелька
	logrus.Infof("Mnemonic: %s", mnemonic.Words())
	logrus.Infof("Wallet Address: %s", account.Address())

	// Создаем новый объект голосования с данными о кошельке голосования
	voteWithID := models.Vote{
		Title:         vote.Title,
		Subtitle:      vote.Subtitle,
		Description:   vote.Description,
		Voter:         vote.Voter,
		Choice:        vote.Choice,
		WalletAddress: account.Address(),
	}

	votePower, err := services.GetVoteStrength(vote.Voter)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	voteWithID.VotePower = votePower

	id, err := services.CreateVote(voteWithID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	voteWithID.ID = id

	utils.JSONResponse(c, http.StatusCreated, voteWithID)
}

// GetVoteHandler получает голосование по ID
func GetVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoteID"})
		return
	}

	vote, err := services.GetVote(id)
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, vote)
}

// DeleteVoteHandler обрабатывает DELETE /votes/:id запрос для удаления голосования
func DeleteVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoteID"})
		return
	}

	if err := services.DeleteVote(id); err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusNoContent, gin.H{})
}

// AddUserVoteHandler добавляет голос пользователя к голосованию
func AddUserVoteHandler(c *gin.Context) {
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid VoteID"})
		return
	}

	var userVote models.UserVote
	if err := c.BindJSON(&userVote); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userVote.VoteID = voteID
	votePower, err := services.GetVoteStrength(userVote.Voter)
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Error determining vote strength"})
		return
	}
	userVote.VotePower = votePower

	id, err := services.AddUserVote(userVote)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to add vote"})
		return
	}
	
	userVote.VoterID = id
	utils.JSONResponse(c, http.StatusCreated, userVote)
}

// GetUserVotesHandler обрабатывает GET /votes/:id/votes запрос для получения всех голосов пользователей для голосования
func GetUserVotesHandler(c *gin.Context) {
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid VoteID"})
		return
	}

	userVotes, err := services.GetUserVotes(voteID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to get user votes"})
		return
	}

	utils.JSONResponse(c, http.StatusOK, userVotes)
}
