package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"payment-service/service"
)

type MidtransController interface {
	PaymentHandler(c *gin.Context)
}

type MidtransControllerImpl struct {
	service.MidtransService
	*amqp.Channel
}

func NewMidtransController(midtransService service.MidtransService, channel *amqp.Channel) *MidtransControllerImpl {
	return &MidtransControllerImpl{MidtransService: midtransService, Channel: channel}
}

func (controller *MidtransControllerImpl) PaymentHandler(c *gin.Context) {
	var notificationPayload map[string]any
	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{"Message": "Bad Request"})
	}
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		c.JSON(http.StatusBadRequest, map[string]any{"Message": "Bad Request"})
	}
	response, _ := controller.MidtransService.VerifyPayment(orderId)
	if response != "" {
		err := controller.Channel.PublishWithContext(context.Background(), "payment", "payment-response", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte(response)})
		if err != nil {
			log.Fatal(err.Error())
		}
		c.JSON(http.StatusOK, map[string]any{"Message": "OK"})
	}
}
