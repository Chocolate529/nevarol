package email

import (
"fmt"
"log"
"net/smtp"
"os"
"strings"
)

// Config holds email configuration
type Config struct {
SMTPHost     string
SMTPPort     string
SMTPUser     string
SMTPPassword string
FromEmail    string
FromName     string
ToEmail      string // Admin email to receive order notifications
}

// NewConfig creates email configuration from environment variables
func NewConfig() *Config {
return &Config{
SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
SMTPPort:     getEnv("SMTP_PORT", "587"),
SMTPUser:     getEnv("SMTP_USER", ""),
SMTPPassword: getEnv("SMTP_PASSWORD", ""),
FromEmail:    getEnv("FROM_EMAIL", ""),
FromName:     getEnv("FROM_NAME", "Transpalet Wheels"),
ToEmail:      getEnv("ADMIN_EMAIL", ""), // Admin email for order notifications
}
}

func getEnv(key, defaultValue string) string {
if value := os.Getenv(key); value != "" {
return value
}
return defaultValue
}

// IsConfigured checks if email is properly configured
func (c *Config) IsConfigured() bool {
return c.SMTPUser != "" && c.SMTPPassword != "" && c.FromEmail != "" && c.ToEmail != ""
}

// OrderDetails contains information about an order for email
type OrderDetails struct {
OrderID       int
CustomerEmail string
CustomerName  string
Phone         string
Address       string
TotalPrice    float64
Items         []OrderItemDetail
}

// OrderItemDetail represents a single item in the order email
type OrderItemDetail struct {
ProductName string
Quantity    int
Price       float64
}

// SendOrderNotification sends order notification email to admin
func (c *Config) SendOrderNotification(order OrderDetails) error {
if !c.IsConfigured() {
log.Println("Email not configured - skipping email notification")
return nil
}

subject := fmt.Sprintf("New Order #%d from %s", order.OrderID, order.CustomerEmail)

// Build item list
var itemsList strings.Builder
for _, item := range order.Items {
itemsList.WriteString(fmt.Sprintf("- %s x%d @ €%.2f = €%.2f\n", 
item.ProductName, item.Quantity, item.Price, item.Price*float64(item.Quantity)))
}

body := fmt.Sprintf(`New Order Received!

Order ID: #%d
Customer Email: %s
Customer Name: %s
Phone: %s
Shipping Address: %s

Items Ordered:
%s
Total: €%.2f

Status: Pending

Please contact the customer to arrange delivery and payment.

---
Transpalet Wheels Order System
`, order.OrderID, order.CustomerEmail, order.CustomerName, order.Phone, 
   order.Address, itemsList.String(), order.TotalPrice)

return c.sendEmail(c.ToEmail, subject, body)
}

// SendOrderConfirmation sends order confirmation email to customer
func (c *Config) SendOrderConfirmation(order OrderDetails) error {
if !c.IsConfigured() {
log.Println("Email not configured - skipping customer confirmation")
return nil
}

subject := fmt.Sprintf("Order Confirmation #%d - Transpalet Wheels", order.OrderID)

// Build item list
var itemsList strings.Builder
for _, item := range order.Items {
itemsList.WriteString(fmt.Sprintf("- %s x%d @ €%.2f = €%.2f\n", 
item.ProductName, item.Quantity, item.Price, item.Price*float64(item.Quantity)))
}

body := fmt.Sprintf(`Thank you for your order!

Order ID: #%d

We have received your order and will contact you shortly to arrange delivery and payment.

Order Details:
%s
Total: €%.2f

Your Contact Information:
Name: %s
Email: %s
Phone: %s
Shipping Address: %s

We will be in touch soon to finalize the details.

Best regards,
Transpalet Wheels Team
`, order.OrderID, itemsList.String(), order.TotalPrice,
   order.CustomerName, order.CustomerEmail, order.Phone, order.Address)

return c.sendEmail(order.CustomerEmail, subject, body)
}

func (c *Config) sendEmail(to, subject, body string) error {
from := c.FromEmail

// Set up authentication
auth := smtp.PlainAuth("", c.SMTPUser, c.SMTPPassword, c.SMTPHost)

// Compose message
msg := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
"To: %s\r\n"+
"Subject: %s\r\n"+
"\r\n"+
"%s\r\n", c.FromName, from, to, subject, body))

// Send email
addr := fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)
err := smtp.SendMail(addr, auth, from, []string{to}, msg)

if err != nil {
log.Printf("Failed to send email to %s: %v", to, err)
return err
}

log.Printf("Email sent successfully to %s", to)
return nil
}
