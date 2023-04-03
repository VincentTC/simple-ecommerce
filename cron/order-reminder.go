package cron

import (
	"context"
	"fmt"
	"log"
)

type cronOrderReminder struct {
	*Cron
}

func orderReminder(c *Cron) *cronOrderReminder {
	return &cronOrderReminder{c}
}

func (c *cronOrderReminder) Handle() {
	ctx := context.Background()
	log.Printf("Cron order reminder start")

	pendingOrders, err := c.options.Sv.GetAllPendingOrdersHandler(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	for _, pendingOrder := range pendingOrders {
		var productNames string
		for i, orderProduct := range pendingOrder.OrderProducts {
			if i > 0 {
				productNames += ", "
			}
			productNames += orderProduct.Product.Name
		}

		link := fmt.Sprintf("ecom://checkout/order/%d", pendingOrder.Id)

		message := fmt.Sprintf("You have pending order for product %s. You can continue to purchase your order on this link %s", productNames, link)

		// TODO: change log print to send email
		log.Println(message)
	}

	log.Printf("Cron order reminder finish")

}
