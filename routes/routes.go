package routes

import (
	"template/handlers"

	"github.com/gofiber/fiber/v2"
)

// Set Route path with functions
func InitRoutes(app *fiber.App) {

	//Without Session Open
	app.Get("/login", handlers.LoginView)
	app.Post("/loginPost", handlers.LoginPost)
	app.Get("/registration", handlers.RegistrationView)
	app.Post("/registrationPost", handlers.RegistrationPost)
	app.Get("/logout", handlers.LogOut)
	app.Get("/reset-password", handlers.ResetPasswordView)
	app.Post("/resetPasswordPost", handlers.ResetPasswordPost)
	//with session
	app.Get("/", handlers.IndexView)
	app.Get("/merchant-transactions", handlers.TransactionsViewNew) //TransactionsView
	app.Get("/withdraw", handlers.WithdrawView)
	app.Post("/withdrawForm", handlers.WithdrawFormPost)
	app.Get("/pdf", handlers.TransPDF)
	app.Get("/excel-download", handlers.TransExcel)
	app.Get("/requested-payment", handlers.RequestedPaymentViews)
	app.Get("/payment-request", handlers.PaymentViews)
	app.Post("/paymentRequestPost", handlers.PaymentRequestPost)
	app.Get("/pay-links", handlers.PayLinksListViews)
	app.Get("/pay-links/new", handlers.PayLinksViews)
	app.Post("/payLinkPost", handlers.PayLinkPost)

	app.Get("/support-tickets", handlers.SupportTicket)
	app.Get("/add-support-ticket", handlers.AddSupportTicket)
	app.Post("/post-support-ticket", handlers.SubmitSupportTicket)
	app.Get("/profile", handlers.ProfileView)
	app.Post("/profilePost", handlers.ProfilePost)
	app.Get("/logged-history", handlers.LoginHistoryView)
	app.Get("/customer", handlers.CustomerView)
	app.Get("/api-key", handlers.ApiKeyView)
	app.Get("/get-api", handlers.GetApiKey)
	app.Get("/change-password", handlers.ChangePasswordView)
	app.Post("/changePasswordPost", handlers.ChangePasswordPost)
	app.Get("/crypto-wallet", handlers.GetCryptoWalletList)
	app.Get("/crypto-wallet/edit/:TID", handlers.EditCryptoWallet)
	app.Get("/crypto-wallet/delete/:TID", handlers.DeleteCryptoWallet)
	app.Get("/add-crypto-wallet", handlers.AddCryptoWalletView)
	app.Post("/cryptoWalletPost", handlers.CryptoWalletPost)
	app.Get("/get-network", handlers.GetNetwork)
	app.Get("/payment", handlers.PaymentQRView)
	app.Post("/pay-data", handlers.PayDataPost)
	app.Get("/pay", handlers.PayQRView)

	app.Get("/success/:TransID", handlers.SuccessView)
	app.Get("/failed/:TransID", handlers.FailedView)
	app.Post("/check-payment-status", handlers.FetchPaymentStatus)
	app.Post("/API/PayLink", handlers.ApiPaymentLink)
	app.Get("/API/CryptoBalance", handlers.ApiBalanceByCrypto)
	app.Get("/API/Transaction/:TransID", handlers.ApiTransactionByTransID)
	app.Get("/API/Transaction", handlers.ApiTransactionList)

	// app.Get("/add-crypto-np", handlers.AddCryptoView)
	// app.Post("/add-crypto-post-np", handlers.AddCryptoPost)
	// app.Get("/request-crypto-np", handlers.RequestCryptoView)
	// app.Post("/request-crypto-post-np", handlers.RequestCryptoPost)
	// app.Get("/transactions-np", handlers.TransactionsNPView)
	// app.Get("/pdf-np", handlers.TransNPPdf)

	/////////////////////////////////////////////////
	app.Get("/admin/login", handlers.AdminLoginView)
	app.Post("/admin/adminLoginPost", handlers.AdminLoginPost)
	app.Get("/admin/logout", handlers.AdminLogOut)
	app.Get("/admin/invoice-list", handlers.AdminInvoiceListView)
	app.Get("/admin/transactions", handlers.AdminTransactionsView)
	app.Get("/admin/trans-details", handlers.AdminTransDetailsView)
	app.Get("/admin/withdraw-list", handlers.AdminWithdrawsView)
	app.Get("/admin/revenue", handlers.AdminRevenueView)
	app.Post("/admin/transApprovalPost", handlers.AdminTransApprovalPost)
	app.Get("/admin/members", handlers.AdminMembersView)
	app.Get("/admin/members-details/:MID", handlers.AdminMembersDetailsView)
	app.Post("/admin/FeesPost", handlers.FeesPost)
	app.Get("/admin/member-status", handlers.AdminMemberStatus)
	app.Get("/admin/login-history", handlers.AdminLoginHistory)
	app.Get("/admin/reset-password", handlers.AdminResetPasswordView)
	app.Post("/admin/resetPasswordPost", handlers.AdminResetPasswordPost)
	// // For Crypto Coin List
	app.Get("/admin/coin-list/edit/:TID", handlers.EditCoin)
	app.Get("/admin/coin-list/delete/:TID", handlers.DeleteCoin)
	app.Get("/admin/coin-list", handlers.GetCoinList)
	app.Get("/admin/add-coin", handlers.AddCoinView)
	app.Post("/admin/coinPost", handlers.CoinPost)
	// // For Crypto Address
	app.Get("/admin/manage-address/edit/:TID", handlers.EditAddress)
	app.Get("/admin/manage-address/delete/:TID", handlers.DeleteAddress)
	app.Get("/admin/manage-address", handlers.GetAddressList)
	app.Get("/admin/add-address", handlers.AddAddressView)
	app.Post("/admin/addressPost", handlers.AddressPost)
	// // For Currency List
	app.Get("/admin/currency/edit/:TID", handlers.EditCurrency)
	app.Get("/admin/currency/delete/:TID", handlers.DeleteCurrency)
	app.Get("/admin/currency", handlers.GetCurrencyList)
	app.Get("/admin/add-currency", handlers.AddCurrencyView)
	app.Post("/admin/currencyPost", handlers.CurrencyPost)
	// For Crypro Currency List
	app.Get("/admin/crypto-currency/edit/:TID", handlers.EditCryptoCurrency)
	app.Get("/admin/crypto-currency/delete/:TID", handlers.DeleteCryptoCurrency)
	app.Get("/admin/crypto-currency", handlers.GetCryptoCurrencyList)
	app.Get("/admin/add-crypto-currency", handlers.AddCryptoCurrencyView)
	app.Post("/admin/cryptoCurrencyPost", handlers.CryptoCurrencyPost)
	// For Admin Manager
	app.Get("/admin/add-admin", handlers.AddAdminView)
	app.Post("/admin/adminPost", handlers.AdminPost)
	app.Get("/admin/admin-manager/edit/:TID", handlers.EditAdmin)
	app.Get("/admin/admin-manager/delete/:TID", handlers.DeleteAdmin)
	app.Get("/admin/admin-manager", handlers.GetAdminList)
	// For Email Template

	app.Get("/admin/add-email-template", handlers.AddTemplateView)
	app.Post("/admin/emailTemplatePost", handlers.TemplatePost)
	app.Get("/admin/email-template/edit/:TID", handlers.EditTemplate)
	app.Get("/admin/email-template/delete/:TID", handlers.DeleteTemplate)
	app.Get("/admin/email-template", handlers.GetTemplateList)
	// For Admin Profile
	app.Get("/admin/profile", handlers.ProfileForm)
	app.Post("/admin/profilePost", handlers.ProfileFormPost)
	// For Admin Change password
	app.Get("/admin/change-password", handlers.ChangePasswordForm)
	app.Post("/admin/change-password-Post", handlers.ChangePasswordFormPost)
	// For Admin Support Ticket
	app.Get("/admin/support-ticket", handlers.SupportTicketListing)
	// For Admin Acquirer
	app.Get("/admin/add-acquirer", handlers.AddAcquirerForm)
	app.Post("/admin/acquirerPost", handlers.AcquirerPost)
	app.Get("/admin/acquirer/edit/:TID", handlers.EditAcquirer)
	app.Get("/admin/acquirer/delete/:TID", handlers.DeleteAcquirer)
	app.Get("/admin/acquirer", handlers.AcquirerList)
	app.Get("/admin/fireblocks-users", handlers.UsersView)
	app.Get("/admin/vault-accounts", handlers.VaultAccountsView)
	app.Get("/admin/vault", handlers.VoltView)
	app.Get("/admin/generate-new-wallet-address/:VID/:WID", handlers.CreateVaultWalletAddress)
	app.Get("/admin/generate-new-wallet-address/:VID/", handlers.CreateVaultWalletView)
	app.Post("/admin/generate-new-wallet-address", handlers.CreateVaultWallet)
	app.Get("/admin/wallet/:VID/:WID", handlers.WalletView)
	// //END Fireblocks API
	app.Get("/admin/", handlers.AdminIndexView)
}
