UPDATE site_level_users slu SET slu.Modules = REPLACE(slu.Modules, 'upi', 'upiSale,upiVoid,upiRefund') WHERE slu.Modules REGEXP 'upi(,|$)';
UPDATE tid_user_override tuo SET tuo.Modules = REPLACE(tuo.Modules, 'upi', 'upiSale,upiVoid,upiRefund') WHERE tuo.Modules REGEXP 'upi(,|$)';