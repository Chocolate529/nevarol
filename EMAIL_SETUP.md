# Email Order Notification Setup Guide

## Overview

The application now sends email notifications when orders are placed:
- **Admin notification**: You receive an email with order details and customer contact info
- **Customer confirmation**: Customer receives order confirmation with their contact details

**Important**: Email is **optional**. Orders will work without email configuration, but notifications won't be sent.

## Quick Setup (Gmail)

### Step 1: Enable 2-Factor Authentication
1. Go to your Google Account: https://myaccount.google.com
2. Select Security
3. Enable 2-Step Verification

### Step 2: Generate App Password
1. Visit: https://myaccount.google.com/apppasswords
2. Select "Mail" and "Other (Custom name)"
3. Enter "Nevarol App" as the name
4. Click "Generate"
5. Copy the 16-character password

### Step 3: Configure Environment Variables

Edit your `.env` file (or set environment variables):

```env
# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=abcd efgh ijkl mnop  # The app password from step 2
FROM_EMAIL=your-email@gmail.com
FROM_NAME=Transpalet Wheels
ADMIN_EMAIL=your-admin-email@gmail.com  # Where order notifications are sent
```

### Step 4: Restart the Application

```bash
# If using Docker
docker-compose restart app

# If running manually
./app
```

## What Emails Are Sent?

### 1. Admin Notification Email
Sent to `ADMIN_EMAIL` when an order is placed:

```
Subject: New Order #123 from customer@email.com

New Order Received!

Order ID: #123
Customer Email: customer@email.com
Customer Name: John Doe
Phone: +1234567890
Shipping Address: 123 Main St, City, Country

Items Ordered:
- Polyurethane Wheel Ø80mm x2 @ €19.99 = €39.98
- Nylon Wheel Ø70mm x1 @ €14.50 = €14.50

Total: €54.48

Status: Pending

Please contact the customer to arrange delivery and payment.
```

### 2. Customer Confirmation Email
Sent to the customer's email address:

```
Subject: Order Confirmation #123 - Transpalet Wheels

Thank you for your order!

Order ID: #123

We have received your order and will contact you shortly to arrange delivery and payment.

Order Details:
- Polyurethane Wheel Ø80mm x2 @ €19.99 = €39.98
- Nylon Wheel Ø70mm x1 @ €14.50 = €14.50

Total: €54.48

Your Contact Information:
Name: John Doe
Email: customer@email.com
Phone: +1234567890
Shipping Address: 123 Main St, City, Country

We will be in touch soon to finalize the details.

Best regards,
Transpalet Wheels Team
```

## Using Other Email Providers

### SendGrid
```env
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your-sendgrid-api-key
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Transpalet Wheels
ADMIN_EMAIL=orders@yourdomain.com
```

### Mailgun
```env
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@yourdomain.com
SMTP_PASSWORD=your-mailgun-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Transpalet Wheels
ADMIN_EMAIL=orders@yourdomain.com
```

### AWS SES
```env
SMTP_HOST=email-smtp.us-east-1.amazonaws.com
SMTP_PORT=587
SMTP_USER=your-ses-smtp-username
SMTP_PASSWORD=your-ses-smtp-password
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Transpalet Wheels
ADMIN_EMAIL=orders@yourdomain.com
```

## Testing Email Configuration

### Method 1: Check Application Logs
When the app starts, you'll see:
```
INFO	Email notification system configured
```

Or if not configured:
```
INFO	Email not configured - orders will be created without email notifications
```

### Method 2: Place a Test Order
1. Start the application
2. Add items to cart
3. Click checkout
4. Fill in contact information:
   - Name: Test Customer
   - Email: your-test-email@gmail.com
   - Phone: 1234567890
   - Address: Test Address
5. Place order
6. Check both admin and customer emails

## Troubleshooting

### "Email not configured" message
**Problem**: Some environment variables are missing.

**Solution**: Ensure ALL email variables are set (SMTP_USER, SMTP_PASSWORD, FROM_EMAIL, ADMIN_EMAIL).

### Emails not being sent (no error)
**Problem**: Gmail blocked the login attempt.

**Solution**: 
1. Check "Less secure app access" is enabled (not recommended)
2. OR use App Passwords (recommended - see Quick Setup above)
3. Check your Gmail "Security" page for blocked sign-in attempts

### "Invalid credentials" error
**Problem**: Wrong SMTP username or password.

**Solution**:
1. Verify SMTP_USER matches your email exactly
2. If using Gmail, ensure you're using an App Password, not your regular password
3. Remove any spaces from the app password

### Emails go to spam
**Problem**: Email provider doesn't trust the sender.

**Solution**:
1. Ask recipients to mark as "Not Spam"
2. Set up SPF and DKIM records for your domain
3. Use a dedicated email service like SendGrid or Mailgun

### Connection timeout
**Problem**: Port 587 might be blocked.

**Solution**:
1. Try port 465 (SSL) instead of 587 (TLS)
2. Update SMTP_PORT=465
3. Check firewall settings

## Running Without Email

Email is completely optional. If you don't configure email:
- Orders will still be created successfully
- Order details are saved in the database
- You can view orders in the admin panel / database
- No emails will be sent (no errors)

This allows you to:
1. Launch the app immediately without email setup
2. Manually process orders by checking the database
3. Add email later when ready

## Security Notes

1. **Never commit** your `.env` file to git (it's in `.gitignore`)
2. **Use App Passwords** for Gmail, not your account password
3. **Use environment variables** in production, not `.env` files
4. **Rotate passwords** regularly
5. **Monitor** email sending logs for suspicious activity

## Support

If you continue to have issues:
1. Check application logs: `docker-compose logs app` or `./app`
2. Verify environment variables are loaded: check startup messages
3. Test SMTP connection using a tool like `telnet smtp.gmail.com 587`
