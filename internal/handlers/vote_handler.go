package handlers

import (
	"bitbucket.org/decimalteam/decimal-go-sdk/wallet"
	"dao_vote/internal/models"
	"dao_vote/internal/services"
	"dao_vote/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	// Генерация нового кошелька с использованием SDK
	account, err := wallet.NewAccount("")
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Не удалось создать кошелек"})
		return
	}

	// Создаем новый объект голосования с данными о кошельке голосования
	voteWithID := models.Vote{
		Title:          vote.Title,
		Subtitle:       vote.Subtitle,
		Description:    vote.Description,
		Voter:          vote.Voter,
		Choice:         vote.Choice,
		WalletMnemonic: account.Mnemonic().Words(),   // Сохраняю сид фразу
		WalletAddress:  account.Address(),            // Сохраняю адрес
		PublicKey:      account.PublicKey().String(), // Сохраняю публичный ключ
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

//

// GetVoteHandler Получает голосование по VoterID
func GetVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoterID"})
		return
	}

	vote, err := services.GetVote(id)
	if err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, vote)
}

// DeleteVoteHandler Создает запрос для удаления голосования от создателя
func DeleteVoteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "invalid VoterID"})
		return
	}

	if err := services.DeleteVote(id); err != nil {
		utils.JSONResponse(c, http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusNoContent, gin.H{})
}

// AddUserVoteHandler Добавляет голоса пользователя к голосованию
func AddUserVoteHandler(c *gin.Context) {
	voteID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid vote VoterID"})
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
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid vote VoterID"})
		return
	}

	userVotes, err := services.GetUserVotes(voteID)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to get user votes"})
		return
	}

	utils.JSONResponse(c, http.StatusOK, userVotes)
}
