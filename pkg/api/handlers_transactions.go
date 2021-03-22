package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gustavooferreira/pgw-payment-processor-service/pkg/core"
)

// AuthoriseTransaction handles authorisation of transactions.
func (s *Server) AuthoriseTransaction(c *gin.Context) {
	requestBody := struct {
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

	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	responseBody := struct {
		Code            uint   `json:"code"`
		AuthorisationID string `json:"authorisation_id,omitempty"`
	}{}

	// Check if we should fail
	if ok := s.Repo.ShouldFail(requestBody.CreditCard.Number, core.CCFailReason_Authorise); ok {
		responseBody.Code = 2
	} else {
		uid := s.Authoriser.Authorise(requestBody.CreditCard.Number)
		responseBody.Code = 1
		responseBody.AuthorisationID = uid
	}

	c.JSON(200, responseBody)
}

// CaptureTransaction handles capturing of transactions.
func (s *Server) CaptureTransaction(c *gin.Context) {
	requestBody := struct {
		AuthorisationID string  `json:"authorisation_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	responseBody := struct {
		Code uint `json:"code"`
	}{}

	responseBody.Code = 1

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(requestBody.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Capture); ok {
			responseBody.Code = 2
		}
	}

	c.JSON(200, responseBody)
}

// VoidTransaction handles voiding transactions.
func (s *Server) VoidTransaction(c *gin.Context) {
	requestBody := struct {
		AuthorisationID string `json:"authorisation_id" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	responseBody := struct {
		Code uint `json:"code"`
	}{}

	responseBody.Code = 1

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(requestBody.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Void); ok {
			responseBody.Code = 2
		}
	}

	c.JSON(200, responseBody)
}

// RefundTransaction handles refunding of transactions.
func (s *Server) RefundTransaction(c *gin.Context) {
	requestBody := struct {
		AuthorisationID string  `json:"authorisation_id" binding:"required"`
		Amount          float64 `json:"amount" binding:"required"`
	}{}

	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		s.Logger.Info(fmt.Sprintf("error parsing body: %s", err.Error()))
		RespondWithError(c, 400, "error parsing body")
		return
	}

	responseBody := struct {
		Code uint `json:"code"`
	}{}

	responseBody.Code = 1

	// Check if we should fail
	if ccNumber, ok := s.Authoriser.GetAssociatedCreditCard(requestBody.AuthorisationID); ok {
		if ok := s.Repo.ShouldFail(ccNumber, core.CCFailReason_Refund); ok {
			responseBody.Code = 2
		}
	}

	c.JSON(200, responseBody)
}
