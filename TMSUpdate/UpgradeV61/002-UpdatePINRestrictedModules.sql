UPDATE data_element SET options = "sale|refund|void|preAuthSale|preAuthCompletion|preAuthCancel|gratuitySale|gratuityCompletion|alipaySale|alipayVoid|alipayRefund|upiSale|upiVoid|upiRefund|eppVoid|X-Read|Z-Read|visaQrSale|mastercardQrSale" WHERE name = 'PINRestrictedModules';