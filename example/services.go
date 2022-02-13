package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mobilepay"
	"net/http"
	"os"
	"strconv"
	"time"
)

// TODO: move this out of the examples
const mobilePayDateQueryFormat = "2006-01-02T15:04"

type PaymentsResource struct {
	mp *mobilepay.Client
}

func (pr PaymentsResource) Routes(r *gin.Engine) {
	router := r.Group("/payments")

	router.GET("/", pr.List)    // GET /payments - Read a list of payments.
	router.POST("/", pr.Create) // POST /payments - Create a new payment

	payment := router.Group("/:id")
	{
		payment.GET("/", pr.GetPayment)           // GET /payments/{id} - Read a single payment by :id.
		payment.GET("/cancel", pr.Cancel)         // GET /payments/{id}/cancel - Cancel a single payment by :id.
		payment.POST("/capture", pr.Capture)      // POST /payments/{id}/capture - Capture a single payment by :id.
		payment.GET("/refunds", pr.Refunds)       // GET /payments/{id}/refunds - Get a list of payment refunds.
		payment.POST("/refunds", pr.CreateRefund) // POST /payments/{id}/refunds - Create a single payment refund by :id.
	}
}

func (pr *PaymentsResource) List(c *gin.Context) {
	pageSize := 100
	perPage := c.Query("per_page")
	if _perPage, err := strconv.Atoi(perPage); err == nil {
		pageSize = _perPage
	}

	pageNumber := 1
	page := c.Query("page")
	if _pageNumber, err := strconv.Atoi(page); err == nil {
		pageNumber = _pageNumber
	}

	opts := mobilepay.ListOptions{
		PageSize:   pageSize,
		PageNumber: pageNumber,
	}

	list, err := pr.mp.Payment.List(c, opts)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list.Payments})
}

func (pr *PaymentsResource) Create(c *gin.Context) {
	type CreatePaymentInput struct {
		Amount      int    `json:"amount" binding:"required"`
		RedirectUri string `json:"redirectUri" binding:"required"`
		Reference   string `json:"reference" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	var input CreatePaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentPointId := os.Getenv("MOBILEPAY_PAYMENT_POINT_ID")

	data := &mobilepay.PaymentParams{
		Amount:         input.Amount,
		IdempotencyKey: c.GetHeader("idempotency-key"),
		PaymentPointId: paymentPointId,
		RedirectUri:    input.RedirectUri,
		Reference:      input.Reference,
		Description:    input.Description,
	}

	payment, err := pr.mp.Payment.Create(c, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (pr *PaymentsResource) GetPayment(c *gin.Context) {
	paymentId := c.Param("id")

	payment, err := pr.mp.Payment.Get(c, paymentId)

	if err != nil {
		log.Println(err)

		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payment})
}

func (pr *PaymentsResource) Refunds(c *gin.Context) {
	paymentId := c.Param("id")

	pageSize := 100
	perPage := c.Query("per_page")
	if _perPage, err := strconv.Atoi(perPage); err == nil {
		pageSize = _perPage
	}

	pageNumber := 1
	page := c.Query("page")
	if _pageNumber, err := strconv.Atoi(page); err == nil {
		pageNumber = _pageNumber
	}

	createdBefore := c.Query("createdBefore")
	_createdBefore, err := time.Parse(mobilePayDateQueryFormat, createdBefore)
	if err == nil {
		createdBefore = _createdBefore.String()
	}

	createdAfter := c.Query("createdAfter")
	_createdAfter, err := time.Parse(mobilePayDateQueryFormat, createdAfter)
	if err == nil {
		createdAfter = _createdAfter.String()
	}

	paymentPointId := os.Getenv("MOBILEPAY_PAYMENT_POINT_ID")

	opts := &mobilepay.RefundsListOptions{
		ListOptions: mobilepay.ListOptions{
			PageSize:   pageSize,
			PageNumber: pageNumber,
		},
		PaymentId:      paymentId,
		PaymentPointId: paymentPointId,
		CreatedBefore:  createdBefore,
		CreatedAfter:   createdAfter,
	}
	data, err := pr.mp.Payment.Refund.List(c, opts)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "an error occurred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data.Refunds})
}

func (pr *PaymentsResource) Capture(c *gin.Context) {
	paymentId := c.Param("id")

	type CaptureInput struct {
		Amount int `json:"amount" binding:"required"`
	}

	var input CaptureInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := pr.mp.Payment.Capture(c, paymentId, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not capture payment"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (pr *PaymentsResource) Cancel(c *gin.Context) {
	paymentId := c.Param("id")

	err := pr.mp.Payment.Cancel(c, paymentId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not cancel payment"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (pr *PaymentsResource) CreateRefund(c *gin.Context) {
	paymentId := c.Param("id")

	type RefundInput struct {
		Amount      int    `json:"amount" binding:"required"`
		Reference   string `json:"reference" binding:"required"`
		Description string `json:"description"`
	}

	var input RefundInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data := &mobilepay.RefundParams{
		PaymentId:      paymentId,
		Amount:         input.Amount,
		IdempotencyKey: c.GetHeader("idempotency-key"),
		Reference:      input.Reference,
		Description:    input.Description,
	}

	refund, err := pr.mp.Payment.Refund.Create(c, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not refund payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": refund})
}
