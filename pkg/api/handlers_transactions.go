package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
)

// AuthoriseTransaction handles authorisation of transactions.
func (s *Server) AuthoriseTransaction(c *gin.Context) {
	bodyData := struct {
		CreditCard struct {
			Name        string `json:"name" binding:"required"`
			Number      int64  `json:"number" binding:"required"`
			ExpiryMonth int    `json:"expiry_month" binding:"required"`
			ExpiryYear  int    `json:"expiry_year" binding:"required"`
			CVV         int    `json:"cvv" binding:"required"`
		} `json:"credit_card" binding:"required"`
		Currency string  `json:"currency" binding:"required"`
		Amount   float64 `json:"amount" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	// Check if we should fail
	if ok := s.Repo.ShouldFail(bodyData.CreditCard.Number, core.CCFailReason_Authorise); ok {
		c.JSON(200, gin.H{"code": 2})
		return
	}

	uid := s.Authoriser.Authorise(bodyData.CreditCard.Number)
	c.JSON(200, gin.H{"code": 1, "authorisation_id": uid})
}

// CaptureTransaction handles capturing of transactions.
func (s *Server) CaptureTransaction(c *gin.Context) {
	bodyData := struct {
		AuthorisationID string  `json:"authorisation_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(bodyData.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Capture); ok {
			c.JSON(200, gin.H{"code": 2})
			return
		}
	}

	c.JSON(200, gin.H{"code": 1})
}

// VoidTransaction handles voiding transactions.
func (s *Server) VoidTransaction(c *gin.Context) {
	bodyData := struct {
		AuthorisationID string `json:"authorisation_id" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(bodyData.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Void); ok {
			c.JSON(200, gin.H{"code": 2})
			return
		}
	}

	c.JSON(200, gin.H{"code": 1})
}

// RefundTransaction handles refunding of transactions.
func (s *Server) RefundTransaction(c *gin.Context) {
	bodyData := struct {
		AuthorisationID string  `json:"authorisation_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&bodyData)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(bodyData.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Refund); ok {
			c.JSON(200, gin.H{"code": 2})
			return
		}
	}

	c.JSON(200, gin.H{"code": 1})
}
