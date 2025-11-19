# Quick Hosting Deployment Guide

This guide provides step-by-step deployment instructions for the recommended hosting platforms.

---

## 🚀 Fly.io Deployment (Recommended - $3-5/month)

### Prerequisites
- GitHub account
- Credit card (for Fly.io account, won't be charged on free tier)

### Step 1: Install Fly CLI

**macOS/Linux:**
```bash
curl -L https://fly.io/install.sh | sh
```

**Windows (PowerShell):**
```powershell
iwr https://fly.io/install.ps1 -useb | iex
```

### Step 2: Login to Fly.io

```bash
fly auth login
```

This opens your browser to complete authentication.

### Step 3: Prepare Your Application

In your project directory (`/home/runner/work/nevarol/nevarol`):

```bash
# Make sure you have a Dockerfile (you already do!)
# Make sure you have a .dockerignore file
```

Your existing `Dockerfile` is already configured correctly!

### Step 4: Launch Application

```bash
# Create and configure the app
fly launch --name nevarol-shop

# Follow the prompts:
# - Choose app name: nevarol-shop (or your preferred name)
# - Choose region: Choose closest to your customers
# - Setup PostgreSQL: YES
# - Deploy now: NO (we'll set secrets first)
```

This creates a `fly.toml` configuration file.

### Step 5: Create PostgreSQL Database

```bash
# Create database
fly postgres create --name nevarol-db --region iad

# Attach to your app
fly postgres attach nevarol-db -a nevarol-shop
```

### Step 6: Set Environment Variables

```bash
# Required: Production flag
fly secrets set IN_PRODUCTION=true -a nevarol-shop

# Optional: Email configuration (recommended)
fly secrets set SMTP_USER=your-email@gmail.com -a nevarol-shop
fly secrets set SMTP_PASSWORD=your-app-password -a nevarol-shop
fly secrets set ADMIN_EMAIL=your-email@gmail.com -a nevarol-shop
```

### Step 7: Deploy!

```bash
fly deploy -a nevarol-shop
```

### Step 8: Open Your Site

```bash
fly open -a nevarol-shop
```

### Step 9: Add Custom Domain (Optional)

```bash
# Add your domain
fly certs add yourdomain.com -a nevarol-shop

# Update your DNS:
# Add A record: @ -> [IP from Fly]
# Add AAAA record: @ -> [IPv6 from Fly]
```

### Useful Fly Commands

```bash
# View logs
fly logs -a nevarol-shop

# Check status
fly status -a nevarol-shop

# SSH into container
fly ssh console -a nevarol-shop

# Scale resources
fly scale memory 512 -a nevarol-shop

# View database connection
fly postgres connect -a nevarol-db
```

---

## 🚂 Railway Deployment (Alternative - $5/month)

### Step 1: Create Railway Account

1. Go to [railway.app](https://railway.app)
2. Sign up with GitHub
3. Verify your account

### Step 2: Create New Project

1. Click "New Project"
2. Select "Deploy from GitHub repo"
3. Choose your `nevarol` repository
4. Railway auto-detects Dockerfile

### Step 3: Add PostgreSQL

1. Click "New" → "Database" → "PostgreSQL"
2. Railway automatically sets `DATABASE_URL`

### Step 4: Set Environment Variables

In Railway dashboard:

1. Go to your service
2. Click "Variables" tab
3. Add:
   ```
   IN_PRODUCTION=true
   SMTP_USER=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   ADMIN_EMAIL=your-email@gmail.com
   ```

### Step 5: Deploy

Railway automatically deploys when you push to GitHub!

```bash
git push origin main
```

### Step 6: Get URL

1. Go to "Settings" tab
2. Click "Generate Domain"
3. Your app is live at `*.railway.app`

### Step 7: Custom Domain (Optional)

1. Go to "Settings" → "Domains"
2. Add custom domain
3. Update DNS to point to Railway

---

## 🌊 DigitalOcean App Platform ($20/month)

### Step 1: Create DigitalOcean Account

1. Go to [digitalocean.com](https://digitalocean.com)
2. Sign up (get $200 free credit with promo)

### Step 2: Create App

1. Click "Apps" → "Create App"
2. Connect GitHub repository
3. Select `nevarol` repo
4. Select branch: `main`

### Step 3: Configure App

1. **Detect Dockerfile**: Auto-detected
2. **Select Plan**: Basic ($5/month)
3. **Region**: Choose closest to customers

### Step 4: Add Database

1. Click "Add Database"
2. Select "PostgreSQL"
3. Choose plan ($15/month for production)

### Step 5: Environment Variables

Add in "Environment Variables" section:
```
IN_PRODUCTION=true
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
ADMIN_EMAIL=your-email@gmail.com
```

### Step 6: Create App

Click "Create Resources"

DigitalOcean builds and deploys automatically.

### Step 7: Custom Domain

1. Go to "Settings" → "Domains"
2. Add your domain
3. Update DNS records

---

## 📧 Email Setup (All Platforms)

### Using Gmail (Recommended for Small Sites)

1. **Enable 2-Factor Authentication**
   - Go to [myaccount.google.com](https://myaccount.google.com)
   - Security → 2-Step Verification → Turn On

2. **Create App Password**
   - Go to [myaccount.google.com/apppasswords](https://myaccount.google.com/apppasswords)
   - Select "Mail" and "Other"
   - Name it "Nevarol Website"
   - Copy the 16-character password

3. **Add to Environment Variables**
   ```bash
   SMTP_USER=your-email@gmail.com
   SMTP_PASSWORD=abcd efgh ijkl mnop  # 16-char app password
   ADMIN_EMAIL=your-email@gmail.com
   ```

### Using SendGrid (Better for High Volume)

1. Sign up at [sendgrid.com](https://sendgrid.com) (free 100 emails/day)
2. Create API key
3. Set environment variables:
   ```bash
   SMTP_HOST=smtp.sendgrid.net
   SMTP_PORT=587
   SMTP_USER=apikey
   SMTP_PASSWORD=your-sendgrid-api-key
   ADMIN_EMAIL=your-email@sendgrid.com
   ```

---

## 🔒 SSL/HTTPS Setup

**Good news:** All recommended platforms provide **automatic HTTPS**!

- ✅ Fly.io: Automatic Let's Encrypt certificates
- ✅ Railway: Automatic SSL on all domains
- ✅ DigitalOcean: Automatic SSL with App Platform

No manual configuration needed!

---

## 💰 Cost Comparison

### Fly.io
- **Free Tier**: $5/month credit
- **Your Cost**: ~$3-5/month
- **Includes**: App hosting + PostgreSQL (3GB)

### Railway
- **Free Tier**: $5/month credit
- **Your Cost**: ~$5-10/month
- **Includes**: App hosting + PostgreSQL

### DigitalOcean
- **App**: $5/month
- **Database**: $15/month
- **Total**: $20/month

### Recommended: **Fly.io** for best value

---

## 🐛 Troubleshooting

### Database Connection Issues

**Error:** "Failed to connect to database"

**Solution:**
1. Check database is created and attached
2. Verify connection string in logs
3. Ensure migrations run on startup

**Fly.io specific:**
```bash
# Check database status
fly postgres connect -a nevarol-db

# View connection info
fly postgres db list -a nevarol-db
```

### Build Failures

**Error:** "Docker build failed"

**Solution:**
1. Check Dockerfile syntax
2. Ensure all dependencies in go.mod
3. Check build logs for specific errors

```bash
# Local test
docker build -t nevarol .
docker run -p 8080:8080 nevarol
```

### Email Not Sending

**Error:** "Failed to send email"

**Solutions:**
1. Verify SMTP credentials are correct
2. Check Gmail app password (not regular password)
3. Ensure 2FA is enabled on Gmail
4. Check SMTP port (587 for TLS)

### Application Crashes

**Check logs:**

**Fly.io:**
```bash
fly logs -a nevarol-shop
```

**Railway:**
View logs in dashboard

**DigitalOcean:**
View logs in App Platform console

---

## 📊 Monitoring Your Site

### Uptime Monitoring (Free)

Use **UptimeRobot** (free):
1. Sign up at [uptimerobot.com](https://uptimerobot.com)
2. Add monitor for your URL
3. Get alerts if site goes down

### Performance Monitoring

**Fly.io:**
```bash
fly dashboard metrics -a nevarol-shop
```

**Railway:**
Built-in metrics in dashboard

**DigitalOcean:**
Built-in insights in App Platform

---

## 🚀 Deployment Checklist

Before going live:

- [ ] Environment variables set (IN_PRODUCTION=true)
- [ ] Email configuration tested
- [ ] Database migrations run successfully
- [ ] Products added to database
- [ ] Test order flow end-to-end
- [ ] SSL/HTTPS working (automatic)
- [ ] Custom domain configured (if applicable)
- [ ] Uptime monitoring set up
- [ ] Test from different devices/browsers
- [ ] Backup strategy in place

---

## 🔄 Updating Your Site

### Fly.io

```bash
# Make changes to code
git add .
git commit -m "Update products"

# Deploy
fly deploy -a nevarol-shop
```

### Railway

```bash
# Make changes and push
git add .
git commit -m "Update products"
git push origin main

# Railway auto-deploys!
```

### DigitalOcean

```bash
# Push to GitHub
git push origin main

# App Platform auto-deploys!
```

---

## 🎓 Next Steps

1. **Choose your platform** (recommend Fly.io)
2. **Follow deployment steps** above
3. **Configure email** for order notifications
4. **Add your products** to the database
5. **Test thoroughly** before sharing
6. **Share with customers** and start selling!

For more details, see:
- [ARCHITECTURE_RECOMMENDATIONS.md](ARCHITECTURE_RECOMMENDATIONS.md) - Architecture decisions
- [PRODUCTION_READINESS.md](PRODUCTION_READINESS.md) - Production checklist
- [EMAIL_SETUP.md](EMAIL_SETUP.md) - Email configuration details

---

**You're ready to deploy! Choose a platform and follow the steps above. Good luck! 🚀**
