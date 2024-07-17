--multiline

UPDATE profile_data SET datavalue = 
'[
      {
        "cardName": "MASTER",
        "brandCode": 1,
        "minLength": 12,
        "maxLength": 19,
        "tacDenial": "0050000000",
        "tacOnline": "FC50BCF800",
        "ranges": [
          "2221>2720",
          "51>55"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "dccEnabled": true,
        "dccDisclaimer": "I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>",
        "dccPrintFee": false
      },
      {
        "cardName": "VISA",
        "brandCode": 2,
        "minLength": 16,
        "maxLength": 19,
        "tacDenial": "0050000000",
        "tacOnline": "DC4004F800",
        "ranges": [
          "4"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "dccEnabled": true,
        "dccDisclaimer": "I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>",
        "dccPrintFee": true
      },
      {
        "cardName": "AMEX",
        "brandCode": 3,
        "amexMid": true,
        "minLength": 1,
        "maxLength": 19,
        "tacDenial": "0010000000",
        "tacOnline": "DE00FC9800",
        "ranges": [
          "34",
          "37"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "disableOnlinePin": true
      },
      {
        "cardName": "DINERS",
        "brandCode": 4,
        "minLength": 14,
        "maxLength": 19,
        "tacDenial": "0010000000",
        "tacOnline": "FCE09CF800",
        "ranges": [
          "36",
          "6510"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
      },
      {
        "cardName": "UNIONPAY",
        "brandCode": 11,
        "minLength": 16,
        "maxLength": 19,
        "tacDenial": "0050000000",
        "tacOnline": "D84004F800",
        "ranges": [
          "62"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
      },
      {
        "cardName": "JCB",
        "brandCode": 12,
        "minLength": 16,
        "maxLength": 16,
        "tacDenial": "0050000000",
        "tacOnline": "FC60ACF800",
        "ranges": [
          "3528>3589"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
      },
      {
        "cardName": "MERCURY",
        "brandCode": 0,
        "minLength": 16,
        "maxLength": 19,
        "tacDenial": "0010000000",
        "tacOnline": "FFFFFFFFFF",
        "ranges": [
          "6504"
        ],
        "offlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E",
        "onlineChipData": "50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E"
      }
    ]', updated_at = NOW() WHERE data_element_id = 21;