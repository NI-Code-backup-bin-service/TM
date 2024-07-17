--multiline
UPDATE profile_data
SET datavalue = '[
   {
      "cardName":"MASTER",
      "brandCode":1,
      "minLength":12,
      "maxLength":19,
      "tacDenial":"0050000000",
      "tacOnline":"FC50BCF800",
      "ranges":[
         "2221>2720",
         "51>55"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "dccEnabled":true,
      "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>",
      "dccPrintFee":false
   },
   {
      "cardName":"VISA",
      "brandCode":2,
      "minLength":16,
      "maxLength":19,
      "tacDenial":"0050000000",
      "tacOnline":"DC4004F800",
      "ranges":[
         "4"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "dccEnabled":true,
      "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>",
      "dccPrintFee":true
   },
   {
      "cardName":"AMEX",
      "brandCode":3,
      "amexMid":true,
      "minLength":1,
      "maxLength":19,
      "tacDenial":"0010000000",
      "tacOnline":"DE00FC9800",
      "ranges":[
         "34",
         "37"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "disableOnlinePin":true
   },
   {
      "cardName":"UNIONPAY",
      "brandCode":11,
      "minLength":16,
      "maxLength":19,
      "refundsEnabled":false,
      "tacDenial":"0050000000",
      "tacOnline":"D84004F800",
      "ranges":[
         "601382",
         "601428",
         "602907",
         "602969",
         "603265",
         "603367",
         "603601",
         "603694",
         "603708",
         "620000>621799",
         "621977",
         "622126>626999",
         "626202",
         "627066>627067",
         "628200>628899",
         "629100>629399",
         "6858",
         "69075",
         "90>98",
         "990027",
         "998800>998802"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
   },
   {
      "cardName":"JCB",
      "brandCode":12,
      "minLength":16,
      "maxLength":16,
      "tacDenial":"0050000000",
      "tacOnline":"FC60ACF800",
      "ranges":[
         "3528>3589"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
   },
   {
      "cardName":"MERCURY",
      "brandCode":0,
      "minLength":16,
      "maxLength":19,
      "tacDenial":"0010000000",
      "tacOnline":"FFFFFFFFFF",
      "ranges":[
         "650401",
         "650008",
         "650483",
         "656001",
         "978432>978439"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
   },
   {
      "cardName":"DINERS",
      "brandCode":4,
      "minLength":14,
      "maxLength":19,
      "tacDenial":"0010000000",
      "tacOnline":"FCE09CF800",
      "ranges":[
         "30",
         "36",
         "38",
         "60110",
         "60112>60114",
         "601174",
         "601177>601179",
         "601186>601199",
         "644>649",
         "65"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
   },
   {
      "cardName":"MAESTRO",
      "brandCode":1,
      "minLength":12,
      "maxLength":19,
      "tacDenial":"0050000000",
      "tacOnline":"FC50BCF800",
      "ranges":[
         "500000>509999",
         "560000>699999"
      ],
      "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
      "dccEnabled":true,
      "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>",
      "dccPrintFee":false
   }
]'
where profile_id = (SELECT profile_id from profile where name = 'global')
and data_element_id = (SELECT data_element_id from data_element where name = 'cardDefinitions'
and data_group_id = (SELECT data_group_id from data_group where name = 'dualCurrency'));