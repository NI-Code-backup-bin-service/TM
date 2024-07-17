SET @profile_id = 2735;
-- Add the Data Group
INSERT ignore into profile_data_group (profile_id, data_group_id, version, updated_at, updated_by, created_at, created_by) values (@profile_id, (SELECT data_group_id from data_group where name = 'multiDeviceEmv'), 1, NOW(), 'NISuper', NOW(), 'NISuper');
-- contactKeyConfigs
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'contactKeyConfigs' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'multiDeviceEmv')), '[ "XAC_9F220101 DF2005A000000003 DF2203000003 DF2314D34A6A776011C7E7CE3AEC5F03AD2F8CFC5503CC DF240420491231 DF260420080101 DF2180C696034213D7D8546984579D1D0F0EA519CFF8DEFFC429354CF3A871A6F7183F1228DA5C7470C055387100CB935A712C4E2864DF5D64BA93FE7E63E71F25B1E5F5298575EBE1C63AA617706917911DC2A75AC28B251C7EF40F2365912490B939BCA2124A30A28F54402C34AECA331AB67E1E79B285DD5771B5D9FF79EA630B75", "XAC_9F220107 DF2005A000000003 DF2203000003 DF2314B4BC56CC4E88324932CBC643D6898F6FE593B172 DF240420491231 DF260420080101 DF218190A89F25A56FA6DA258C8CA8B40427D927B4A1EB4D7EA326BBB12F97DED70AE5E4480FC9C5E8A972177110A1CC318D06D2F8F5C4844AC5FA79A4DC470BB11ED635699C17081B90F1B984F12E92C1C529276D8AF8EC7F28492097D8CD5BECEA16FE4088F6CFAB4A1B42328A1B996F9278B0B7E3311CA5EF856C2F888474B83612A82E4E00D0CD4069A6783140433D50725F", "XAC_9F220108 DF2005A000000003 DF2203000003 DF231420D213126955DE205ADC2FD2822BD22DE21CF9A8 DF240420491231 DF260420080101 DF2181B0D9FD6ED75D51D0E30664BD157023EAA1FFA871E4DA65672B863D255E81E137A51DE4F72BCC9E44ACE12127F87E263D3AF9DD9CF35CA4A7B01E907000BA85D24954C2FCA3074825DDD4C0C8F186CB020F683E02F2DEAD3969133F06F7845166ACEB57CA0FC2603445469811D293BFEFBAFAB57631B3DD91E796BF850A25012F1AE38F05AA5C4D6D03B1DC2E568612785938BBC9B3CD3A910C1DA55A5A9218ACE0F7A21287752682F15832A678D6E1ED0B", "XAC_9F220109 DF2005A000000003 DF2203000003 DF23141FF80A40173F52D7D27E0F26A146A1C8CCB29046 DF240420491231 DF260420080101 DF2181F89D912248DE0A4E39C1A7DDE3F6D2588992C1A4095AFBD1824D1BA74847F2BC4926D2EFD904B4B54954CD189A54C5D1179654F8F9B0D2AB5F0357EB642FEDA95D3912C6576945FAB897E7062CAA44A4AA06B8FE6E3DBA18AF6AE3738E30429EE9BE03427C9D64F695FA8CAB4BFE376853EA34AD1D76BFCAD15908C077FFE6DC5521ECEF5D278A96E26F57359FFAEDA19434B937F1AD999DC5C41EB11935B44C18100E857F431A4A5A6BB65114F174C2D7B59FDF237D6BB1DD0916E644D709DED56481477C75D95CDD68254615F7740EC07F330AC5D67BCD75BF23D28A140826C026DBDE971A37CD3EF9B8DF644AC385010501EFC6509D7A41", "XAC_9F220103 DF2005A000000004 DF2203000003 DF23148BB99ADDF7B560110955014505FB6B5F8308CE27 DF240420491231 DF260420080101 DF21609E15214212F6308ACA78B80BD986AC287516846C8D548A9ED0A42E7D997C902C3E122D1B9DC30995F4E25C75DD7EE0A0CE293B8CC02B977278EF256D761194924764942FE714FA02E4D57F282BA3B2B62C9E38EF6517823F2CA831BDDF6D363D", "XAC_9F220104 DF2005A000000004 DF2203000003 DF2314381A035DA58B482EE2AF75F4C3F2CA469BA4AA6C DF240420491231 DF260420080101 DF218190A6DA428387A502D7DDFB7A74D3F412BE762627197B25435B7A81716A700157DDD06F7CC99D6CA28C2470527E2C03616B9C59217357C2674F583B3BA5C7DCF2838692D023E3562420B4615C439CA97C44DC9A249CFCE7B3BFB22F68228C3AF13329AA4A613CF8DD853502373D62E49AB256D2BC17120E54AEDCED6D96A4287ACC5C04677D4A5A320DB8BEE2F775E5FEC5", "XAC_9F220105 DF2005A000000004 DF2203000003 DF2314EBFA0D5D06D8CE702DA3EAE890701D45E274C845 DF240420491231 DF260420080101 DF2181B0B8048ABC30C90D976336543E3FD7091C8FE4800DF820ED55E7E94813ED00555B573FECA3D84AF6131A651D66CFF4284FB13B635EDD0EE40176D8BF04B7FD1C7BACF9AC7327DFAA8AA72D10DB3B8E70B2DDD811CB4196525EA386ACC33C0D9D4575916469C4E4F53E8E1C912CC618CB22DDE7C3568E90022E6BBA770202E4522A2DD623D180E215BD1D1507FE3DC90CA310D27B3EFCCD8F83DE3052CAD1E48938C68D095AAC91B5F37E28BB49EC7ED597", "XAC_9F220106 DF2005A000000004 DF2203000003 DF2314F910A1504D5FFB793D94F3B500765E1ABCAD72D9 DF240420491231 DF260420080101 DF2181F8CB26FC830B43785B2BCE37C81ED334622F9622F4C89AAE641046B2353433883F307FB7C974162DA72F7A4EC75D9D657336865B8D3023D3D645667625C9A07A6B7A137CF0C64198AE38FC238006FB2603F41F4F3BB9DA1347270F2F5D8C606E420958C5F7D50A71DE30142F70DE468889B5E3A08695B938A50FC980393A9CBCE44AD2D64F630BB33AD3F5F5FD495D31F37818C1D94071342E07F1BEC2194F6035BA5DED3936500EB82DFDA6E8AFB655B1EF3D0D7EBF86B66DD9F29F6B1D324FE8B26CE38AB2013DD13F611E7A594D675C4432350EA244CC34F3873CBA06592987A1D7E852ADC22EF5A2EE28132031E48F74037E3B34AB747F", "XAC_9F220103 DF2005A000000025 DF2203000003 DF23148708A3E3BBC1BB0BE73EBD8D19D4E5D20166BF6C DF240420491231 DF260420080101 DF2180B0C2C6E2A6386933CD17C239496BF48C57E389164F2A96BFF133439AE8A77B20498BD4DC6959AB0C2D05D0723AF3668901937B674E5A2FA92DDD5E78EA9D75D79620173CC269B35F463B3D4AAFF2794F92E6C7A3FB95325D8AB95960C3066BE548087BCB6CE12688144A8B4A66228AE4659C634C99E36011584C095082A3A3E3", "XAC_9F22010E DF2005A000000025 DF2203000003 DF2314A7266ABAE64B42A3668851191D49856E17F8FBCD DF240420491231 DF260420080101 DF218190AA94A8C6DAD24F9BA56A27C09B01020819568B81A026BE9FD0A3416CA9A71166ED5084ED91CED47DD457DB7E6CBCD53E560BC5DF48ABC380993B6D549F5196CFA77DFB20A0296188E969A2772E8C4141665F8BB2516BA2C7B5FC91F8DA04E8D512EB0F6411516FB86FC021CE7E969DA94D33937909A53A57F907C40C22009DA7532CB3BE509AE173B39AD6A01BA5BB85", "XAC_9F22010F DF2005A000000025 DF2203000003 DF2314A73472B3AB557493A9BC2179CC8014053B12BAB4 DF240420491231 DF260420080101 DF2181B0C8D5AC27A5E1FB89978C7C6479AF993AB3800EB243996FBB2AE26B67B23AC482C4B746005A51AFA7D2D83E894F591A2357B30F85B85627FF15DA12290F70F05766552BA11AD34B7109FA49DE29DCB0109670875A17EA95549E92347B948AA1F045756DE56B707E3863E59A6CBE99C1272EF65FB66CBB4CFF070F36029DD76218B21242645B51CA752AF37E70BE1A84FF31079DC0048E928883EC4FADD497A719385C2BBBEBC5A66AA5E5655D18034EC5", "XAC_9F220110 DF2005A000000025 DF2203000003 DF2314C729CF2FD262394ABC4CC173506502446AA9B9FD DF240420491231 DF260420080101 DF2181F8CF98DFEDB3D3727965EE7797723355E0751C81D2D3DF4D18EBAB9FB9D49F38C8C4A826B99DC9DEA3F01043D4BF22AC3550E2962A59639B1332156422F788B9C16D40135EFD1BA94147750575E636B6EBC618734C91C1D1BF3EDC2A46A43901668E0FFC136774080E888044F6A1E65DC9AAA8928DACBEB0DB55EA3514686C6A732CEF55EE27CF877F110652694A0E3484C855D882AE191674E25C296205BBB599455176FDD7BBC549F27BA5FE35336F7E29E68D783973199436633C67EE5A680F05160ED12D1665EC83D1997F10FD05BBDBF9433E8F797AEE3E9F02A34228ACE927ABE62B8B9281AD08D3DF5C7379685045D7BA5FCDE58637", "XAC_9F220109 DF2005A000000065 DF2203000003 DF23144410C6D51C2F83ADFD92528FA6E38A32DF048D0A DF240420491231 DF260420080101 DF2180B72A8FEF5B27F2B550398FDCC256F714BAD497FF56094B7408328CB626AA6F0E6A9DF8388EB9887BC930170BCC1213E90FC070D52C8DCD0FF9E10FAD36801FE93FC998A721705091F18BC7C98241CADC15A2B9DA7FB963142C0AB640D5D0135E77EBAE95AF1B4FEFADCF9C012366BDDA0455C1564A68810D7127676D493890BD", "XAC_9F220110 DF2005A000000065 DF2203000003 DF2314C75E5210CBE6E8F0594A0F1911B07418CADB5BAB DF240420491231 DF260420080101 DF21819099B63464EE0B4957E4FD23BF923D12B61469B8FFF8814346B2ED6A780F8988EA9CF0433BC1E655F05EFA66D0C98098F25B659D7A25B8478A36E489760D071F54CDF7416948ED733D816349DA2AADDA227EE45936203CBF628CD033AABA5E5A6E4AE37FBACB4611B4113ED427529C636F6C3304F8ABDD6D9AD660516AE87F7F2DDF1D2FA44C164727E56BBC9BA23C0285", "XAC_9F220112 DF2005A000000065 DF2203000003 DF2314874B379B7F607DC1CAF87A19E400B6A9E25163E8 DF240420491231 DF260420080101 DF2181B0ADF05CD4C5B490B087C3467B0F3043750438848461288BFEFD6198DD576DC3AD7A7CFA07DBA128C247A8EAB30DC3A30B02FCD7F1C8167965463626FEFF8AB1AA61A4B9AEF09EE12B009842A1ABA01ADB4A2B170668781EC92B60F605FD12B2B2A6F1FE734BE510F60DC5D189E401451B62B4E06851EC20EBFF4522AACC2E9CDC89BC5D8CDE5D633CFD77220FF6BBD4A9B441473CC3C6FEFC8D13E57C3DE97E1269FA19F655215B23563ED1D1860D8681", "XAC_9F220114 DF2005A000000065 DF2203000003 DF2314C0D15F6CD957E491DB56DCDD1CA87A03EBE06B7B DF240420491231 DF260420080101 DF2181F8AEED55B9EE00E1ECEB045F61D2DA9A66AB637B43FB5CDBDB22A2FBB25BE061E937E38244EE5132F530144A3F268907D8FD648863F5A96FED7E42089E93457ADC0E1BC89C58A0DB72675FBC47FEE9FF33C16ADE6D341936B06B6A6F5EF6F66A4EDD981DF75DA8399C3053F430ECA342437C23AF423A211AC9F58EAF09B0F837DE9D86C7109DB1646561AA5AF0289AF5514AC64BC2D9D36A179BB8A7971E2BFA03A9E4B847FD3D63524D43A0E8003547B94A8A75E519DF3177D0A60BC0B4BAB1EA59A2CBB4D2D62354E926E9C7D3BE4181E81BA60F8285A896D17DA8C3242481B6C405769A39D547C74ED9FF95A70A796046B5EFF36682DC29", "XAC_9F220101 DF2005A000000152 DF2203000003 DF2314E0C2C1EA411DB24EC3E76A9403F0B7B6F406F398 DF240420491231 DF260420080101 DF21808D1727AB9DC852453193EA0810B110F2A3FD304BE258338AC2650FA2A040FA10301EA53DF18FD9F40F55C44FE0EE7C7223BC649B8F9328925707776CB86F3AC37D1B22300D0083B49350E09ABB4B62A96363B01E4180E158EADDD6878E85A6C9D56509BF68F0400AFFBC441DDCCDAF9163C4AACEB2C3E1EC13699D23CDA9D3AD", "XAC_9F220103 DF2005A000000152 DF2203000003 DF2314CA1E9099327F0B786D8583EC2F27E57189503A57 DF240420491231 DF260420080101 DF218190BF321241BDBF3585FFF2ACB89772EBD18F2C872159EAA4BC179FB03A1B850A1A758FA2C6849F48D4C4FF47E02A575FC13E8EB77AC37135030C5600369B5567D3A7AAF02015115E987E6BE566B4B4CC03A4E2B16CD9051667C2CD0EEF4D76D27A6F745E8BBEB45498ED8C30E2616DB4DBDA4BAF8D71990CDC22A8A387ACB21DD88E2CC27962B31FBD786BBB55F9E0B041", "XAC_9F220104 DF2005A000000152 DF2203000003 DF231417F971CAF6B708E5B9165331FBA91593D0C0BF66 DF240420491231 DF260420080101 DF2181B08EEEC0D6D3857FD558285E49B623B109E6774E06E9476FE1B2FB273685B5A235E955810ADDB5CDCC2CB6E1A97A07089D7FDE0A548BDC622145CA2DE3C73D6B14F284B3DC1FA056FC0FB2818BCD7C852F0C97963169F01483CE1A63F0BF899D412AB67C5BBDC8B4F6FB9ABB57E95125363DBD8F5EBAA9B74ADB93202050341833DEE8E38D28BD175C83A6EA720C262682BEABEA8E955FE67BD9C2EFF7CB9A9F45DD5BDA4A1EEFB148BC44FFF68D9329FD", "XAC_9F220105 DF2005A000000152 DF2203000003 DF231412BCD407B6E627A750FDF629EE8C2C9CC7BA636A DF240420491231 DF260420080101 DF2181F8E1200E9F4428EB71A526D6BB44C957F18F27B20BACE978061CCEF23532DBEBFAF654A149701C14E6A2A7C2ECAC4C92135BE3E9258331DDB0967C3D1D375B996F25B77811CCCC06A153B4CE6990A51A0258EA8437EDBEB701CB1F335993E3F48458BC1194BAD29BF683D5F3ECB984E31B7B9D2F6D947B39DEDE0279EE45B47F2F3D4EEEF93F9261F8F5A571AFBFB569C150370A78F6683D687CB677777B2E7ABEFCFC8F5F93501736997E8310EE0FD87AFAC5DA772BA277F88B44459FCA563555017CD0D66771437F8B6608AA1A665F88D846403E4C41AFEEDB9729C2B2511CFE228B50C1B152B2A60BBF61D8913E086210023A3AA499E423", "XAC_9F220101 DF2005A000000333 DF2203000003 DF2314E881E390675D44C2DD81234DCE29C3F5AB2297A0 DF240420491231 DF260420080101 DF2180BBE9066D2517511D239C7BFA77884144AE20C7372F515147E8CE6537C54C0A6A4D45F8CA4D290870CDA59F1344EF71D17D3F35D92F3F06778D0D511EC2A7DC4FFEADF4FB1253CE37A7B2B5A3741227BEF72524DA7A2B7B1CB426BEE27BC513B0CB11AB99BC1BC61DF5AC6CC4D831D0848788CD74F6D543AD37C5A2B4C5D5A93B", "XAC_9F220102 DF2005A000000333 DF2203000003 DF231403BB335A8549A03B87AB089D006F60852E4B8060 DF240420491231 DF260420080101 DF218190A3767ABD1B6AA69D7F3FBF28C092DE9ED1E658BA5F0909AF7A1CCD907373B7210FDEB16287BA8E78E1529F443976FD27F991EC67D95E5F4E96B127CAB2396A94D6E45CDA44CA4C4867570D6B07542F8D4BF9FF97975DB9891515E66F525D2B3CBEB6D662BFB6C3F338E93B02142BFC44173A3764C56AADD202075B26DC2F9F7D7AE74BD7D00FD05EE430032663D27A57", "XAC_9F220103 DF2005A000000333 DF2203000003 DF231487F0CD7C0E86F38F89A66F8C47071A8B88586F26 DF240420491231 DF260420080101 DF2181B0B0627DEE87864F9C18C13B9A1F025448BF13C58380C91F4CEBA9F9BCB214FF8414E9B59D6ABA10F941C7331768F47B2127907D857FA39AAF8CE02045DD01619D689EE731C551159BE7EB2D51A372FF56B556E5CB2FDE36E23073A44CA215D6C26CA68847B388E39520E0026E62294B557D6470440CA0AEFC9438C923AEC9B2098D6D3A1AF5E8B1DE36F4B53040109D89B77CAFAF70C26C601ABDF59EEC0FDC8A99089140CD2E817E335175B03B7AA33D", "XAC_9F220104 DF2005A000000333 DF2203000003 DF2314F527081CF371DD7E1FD4FA414A665036E0F5E6E5 DF240420491231 DF260420080101 DF2181F8BC853E6B5365E89E7EE9317C94B02D0ABB0DBD91C05A224A2554AA29ED9FCB9D86EB9CCBB322A57811F86188AAC7351C72BD9EF196C5A01ACEF7A4EB0D2AD63D9E6AC2E7836547CB1595C68BCBAFD0F6728760F3A7CA7B97301B7E0220184EFC4F653008D93CE098C0D93B45201096D1ADFF4CF1F9FC02AF759DA27CD6DFD6D789B099F16F378B6100334E63F3D35F3251A5EC78693731F5233519CDB380F5AB8C0F02728E91D469ABD0EAE0D93B1CC66CE127B29C7D77441A49D09FCA5D6D9762FC74C31BB506C8BAE3C79AD6C2578775B95956B5370D1D0519E37906B384736233251E8F09AD79DFBE2C6ABFADAC8E4D8624318C27DAF1", "XAC_9F220103 DF2005A000000529 DF2203000003 DF23143D439C45EA44C0AB82A395A71987E1120CAC7A99 DF240420491231 DF260420080101 DF2181B09774DE509FEFF4E0990D4C51A707D3048CC70DB24BD08D1751257DB69AB1760E6115B093A402222BB3AB8CDEA6F3FE2DDD589E1F32E21C84F42BB3271FC6BEBC48FD1735B9E995A9AFF5543C27B0481B8FCCE9F3BCB7E5352BE732D6174BC3D06A4A515190FB1EAF39E8817AFA7DF1F2538E4CFFCFF9604BFB120F345611E30D4FE73420558E7B85A75997AB89E1AE9E273C93784A2709F8B6649D7179759A07D92B80A87522839C7CA56D56D3D6986D", "XAC_9F220105 DF2005A000000529 DF2203000003 DF23146405D7A73641EE74C1B110AD0D77A1DDE7D2674C DF240420491231 DF260420080101 DF2181F89F551F8A28B03BDF7D751087C04CA31772C6B55062810D1BA5E7B52BEC6D6A0E105E8DFDE0E309815D87CEDFFC2C453937C0F0B01FD096BBA39FA3C204E088CE6CDE18B8530BB858B4F034DA6BDC882702773550B74537FB81DF1E9BF44BB47666042646B5C3114208C4FBA1F02C19CC5755CF8AE29FD72E7215517B9CD23EEAAFA9FDF6A9712768A760E5452DA7124C12F89747D15803AF3B141222C2FFFC0F12F48E56CD2CAB7DCB9343FC01E2CD279334281DD85BEAF449A8FD91E5987AF194538F4BD9D16E8178F7BF5E460663A20CA7D2BCA8B7CFD47929220112DEEC23331D6FE7EAA60E4AF2FAC48A1D8A740775A2AB5957885355", "XAC_9F220106 DF2005A000000529 DF2203000003 DF231481ABD4E30D654F229D76A40BC1DD8DC558D9B6F3 DF240420491231 DF260420080101 DF218190CA7167D074A7EB59BF41377F5D37E71C96F0954A6E722642D7EBD21A3A996194BFFD46F718AA22CF192F0B36072FFEE9A4CF47D6241FA67C431871D333485CE721074131BA18EA8AC900A8B4232AF79274A8F47F1D193FD95161647BA7D804996578BE1AD3657EF8F1331777B1342ADF00B2E1FBFBCDEF0C4EA32A93C14E7284BCFBD9C859589C293CB8CE66C80FAE27", "VF_9F0605A0000000039F220199DF050420291231DF028180AB79FCC9520896967E776E64444E5DCDD6E13611874F3985722520425295EEA4BD0C2781DE7F31CD3D041F565F747306EED62954B17EDABA3A6C5B85A1DE1BEB9A34141AF38FCF8279C9DEA0D5A6710D08DB4124F041945587E20359BAB47B7575AD94262D4B25F264AF33DEDCF28E09615E937DE32EDC03C54445FE7E382777DF040103DF03144ABFFD6B1C51212D05552E431C5B17007D2F5E6DBF010131DF070101", "VF_9F0605A0000000039F220195DF050420291231DF028190BE9E1FA5E9A803852999C4AB432DB28600DCD9DAB76DFAAA47355A0FE37B1508AC6BF38860D3C6C2E5B12A3CAAF2A7005A7241EBAA7771112C74CF9A0634652FBCA0E5980C54A64761EA101A114E0F0B5572ADD57D010B7C9C887E104CA4EE1272DA66D997B9A90B5A6D624AB6C57E73C8F919000EB5F684898EF8C3DBEFB330C62660BED88EA78E909AFF05F6DA627BDF040103DF0314EE1511CEC71020A9B90443B37B1D5F6E703030F6BF010131DF070101", "VF_9F0605A0000000039F220192DF050420291231DF0281B0996AF56F569187D09293C14810450ED8EE3357397B18A2458EFAA92DA3B6DF6514EC060195318FD43BE9B8F0CC669E3F844057CBDDF8BDA191BB64473BC8DC9A730DB8F6B4EDE3924186FFD9B8C7735789C23A36BA0B8AF65372EB57EA5D89E7D14E9C7B6B557460F10885DA16AC923F15AF3758F0F03EBD3C5C2C949CBA306DB44E6A2C076C5F67E281D7EF56785DC4D75945E491F01918800A9E2DC66F60080566CE0DAF8D17EAD46AD8E30A247C9FDF040103DF0314429C954A3859CEF91295F663C963E582ED6EB253BF010131DF070101", "VF_9F0605A0000000039F220107DF050420291231DF028190A89F25A56FA6DA258C8CA8B40427D927B4A1EB4D7EA326BBB12F97DED70AE5E4480FC9C5E8A972177110A1CC318D06D2F8F5C4844AC5FA79A4DC470BB11ED635699C17081B90F1B984F12E92C1C529276D8AF8EC7F28492097D8CD5BECEA16FE4088F6CFAB4A1B42328A1B996F9278B0B7E3311CA5EF856C2F888474B83612A82E4E00D0CD4069A6783140433D50725FDF040103DF0314B4BC56CC4E88324932CBC643D6898F6FE593B172BF010131DF070101", "VF_9F0605A0000000039F220108DF050420291231DF0281B0D9FD6ED75D51D0E30664BD157023EAA1FFA871E4DA65672B863D255E81E137A51DE4F72BCC9E44ACE12127F87E263D3AF9DD9CF35CA4A7B01E907000BA85D24954C2FCA3074825DDD4C0C8F186CB020F683E02F2DEAD3969133F06F7845166ACEB57CA0FC2603445469811D293BFEFBAFAB57631B3DD91E796BF850A25012F1AE38F05AA5C4D6D03B1DC2E568612785938BBC9B3CD3A910C1DA55A5A9218ACE0F7A21287752682F15832A678D6E1ED0BDF040103DF031420D213126955DE205ADC2FD2822BD22DE21CF9A8BF010131DF070101", "VF_9F0605A0000000039F220109DF050420291231DF0281F89D912248DE0A4E39C1A7DDE3F6D2588992C1A4095AFBD1824D1BA74847F2BC4926D2EFD904B4B54954CD189A54C5D1179654F8F9B0D2AB5F0357EB642FEDA95D3912C6576945FAB897E7062CAA44A4AA06B8FE6E3DBA18AF6AE3738E30429EE9BE03427C9D64F695FA8CAB4BFE376853EA34AD1D76BFCAD15908C077FFE6DC5521ECEF5D278A96E26F57359FFAEDA19434B937F1AD999DC5C41EB11935B44C18100E857F431A4A5A6BB65114F174C2D7B59FDF237D6BB1DD0916E644D709DED56481477C75D95CDD68254615F7740EC07F330AC5D67BCD75BF23D28A140826C026DBDE971A37CD3EF9B8DF644AC385010501EFC6509D7A41DF040103DF03141FF80A40173F52D7D27E0F26A146A1C8CCB29046BF010131DF070101", "VF_9F0605A0000000039F220101DF050420291231DF028180C696034213D7D8546984579D1D0F0EA519CFF8DEFFC429354CF3A871A6F7183F1228DA5C7470C055387100CB935A712C4E2864DF5D64BA93FE7E63E71F25B1E5F5298575EBE1C63AA617706917911DC2A75AC28B251C7EF40F2365912490B939BCA2124A30A28F54402C34AECA331AB67E1E79B285DD5771B5D9FF79EA630B75DF040103DF0314D34A6A776011C7E7CE3AEC5F03AD2F8CFC5503CCBF010131DF070101", "VF_9F0605A0000000039F220103DF050420291231DF0270B3E5E667506C47CAAFB12A2633819350846697DD65A796E5CE77C57C626A66F70BB630911612AD2832909B8062291BECA46CD33B66A6F9C9D48CED8B4FC8561C8A1D8FB15862C9EB60178DEA2BE1F82236FFCFF4F3843C272179DCDD384D541053DA6A6A0D3CE48FDC2DC4E3E0EEE15FDF040103DF0314FE70AB3B4D5A1B9924228ADF8027C758483A8B7EBF010131DF070101", "VF_9F0605A0000000049F2201EFDF050420291231DF0281F8A191CB87473F29349B5D60A88B3EAEE0973AA6F1A082F358D849FDDFF9C091F899EDA9792CAF09EF28F5D22404B88A2293EEBBC1949C43BEA4D60CFD879A1539544E09E0F09F60F065B2BF2A13ECC705F3D468B9D33AE77AD9D3F19CA40F23DCF5EB7C04DC8F69EBA565B1EBCB4686CD274785530FF6F6E9EE43AA43FDB02CE00DAEC15C7B8FD6A9B394BABA419D3F6DC85E16569BE8E76989688EFEA2DF22FF7D35C043338DEAA982A02B866DE5328519EBBCD6F03CDD686673847F84DB651AB86C28CF1462562C577B853564A290C8556D818531268D25CC98A4CC6A0BDFFFDA2DCCA3A94C998559E307FDDF915006D9A987B07DDAEB3BDF040103DF031421766EBB0EE122AFB65D7845B73DB46BAB65427ABF010131DF070101", "VF_9F0605A0000000049F2201FADF050420291231DF028190A90FCD55AA2D5D9963E35ED0F440177699832F49C6BAB15CDAE5794BE93F934D4462D5D12762E48C38BA83D8445DEAA74195A301A102B2F114EADA0D180EE5E7A5C73E0C4E11F67A43DDAB5D55683B1474CC0627F44B8D3088A492FFAADAD4F42422D0E7013536C3C49AD3D0FAE96459B0F6B1B6056538A3D6D44640F94467B108867DEC40FAAECD740C00E2B7A8852DDF040103DF03142CFBB82409ED86A31973B0E0CEEA381BC43C8097BF010131DF070101", "VF_9F0605A0000000049F220104DF050420291231DF028190A6DA428387A502D7DDFB7A74D3F412BE762627197B25435B7A81716A700157DDD06F7CC99D6CA28C2470527E2C03616B9C59217357C2674F583B3BA5C7DCF2838692D023E3562420B4615C439CA97C44DC9A249CFCE7B3BFB22F68228C3AF13329AA4A613CF8DD853502373D62E49AB256D2BC17120E54AEDCED6D96A4287ACC5C04677D4A5A320DB8BEE2F775E5FEC5DF040103DF0314381A035DA58B482EE2AF75F4C3F2CA469BA4AA6CBF010131DF070101", "VF_9F0605A0000000049F220105DF050420291231DF0281B0B8048ABC30C90D976336543E3FD7091C8FE4800DF820ED55E7E94813ED00555B573FECA3D84AF6131A651D66CFF4284FB13B635EDD0EE40176D8BF04B7FD1C7BACF9AC7327DFAA8AA72D10DB3B8E70B2DDD811CB4196525EA386ACC33C0D9D4575916469C4E4F53E8E1C912CC618CB22DDE7C3568E90022E6BBA770202E4522A2DD623D180E215BD1D1507FE3DC90CA310D27B3EFCCD8F83DE3052CAD1E48938C68D095AAC91B5F37E28BB49EC7ED597DF040103DF0314EBFA0D5D06D8CE702DA3EAE890701D45E274C845BF010131DF070101", "VF_9F0605A0000000049F220106DF050420291231DF0281F8CB26FC830B43785B2BCE37C81ED334622F9622F4C89AAE641046B2353433883F307FB7C974162DA72F7A4EC75D9D657336865B8D3023D3D645667625C9A07A6B7A137CF0C64198AE38FC238006FB2603F41F4F3BB9DA1347270F2F5D8C606E420958C5F7D50A71DE30142F70DE468889B5E3A08695B938A50FC980393A9CBCE44AD2D64F630BB33AD3F5F5FD495D31F37818C1D94071342E07F1BEC2194F6035BA5DED3936500EB82DFDA6E8AFB655B1EF3D0D7EBF86B66DD9F29F6B1D324FE8B26CE38AB2013DD13F611E7A594D675C4432350EA244CC34F3873CBA06592987A1D7E852ADC22EF5A2EE28132031E48F74037E3B34AB747FDF040103DF0314F910A1504D5FFB793D94F3B500765E1ABCAD72D9BF010131DF070101", "VF_9F0605A0000000049F2201F1DF050420231231DF0281B0A0DCF4BDE19C3546B4B6F0414D174DDE294AABBB828C5A834D73AAE27C99B0B053A90278007239B6459FF0BBCD7B4B9C6C50AC02CE91368DA1BD21AAEADBC65347337D89B68F5C99A09D05BE02DD1F8C5BA20E2F13FB2A27C41D3F85CAD5CF6668E75851EC66EDBF98851FD4E42C44C1D59F5984703B27D5B9F21B8FA0D93279FBBF69E090642909C9EA27F898959541AA6757F5F624104F6E1D3A9532F2A6E51515AEAD1B43B3D7835088A2FAFA7BE7DF040103DF0314D8E68DA167AB5A85D8C3D55ECB9B0517A1A5B4BBBF010131DF070101", "VF_9F0605A0000000049F220103DF050420291231DF028180C2490747FE17EB0584C88D47B1602704150ADC88C5B998BD59CE043EDEBF0FFEE3093AC7956AD3B6AD4554C6DE19A178D6DA295BE15D5220645E3C8131666FA4BE5B84FE131EA44B039307638B9E74A8C42564F892A64DF1CB15712B736E3374F1BBB6819371602D8970E97B900793C7C2A89A4A1649A59BE680574DD0B60145DF040103DF03145ADDF21D09278661141179CBEFF272EA384B13BBBF010131DF070101", "VF_9F0605A0000000049F220109DF050420291231DF028180C132F436477A59302E885646102D913EC86A95DD5D0A56F625F472B67F52179BC8BD258A7CD43EF1720AC0065519E3FFCECC26F978EDF9FB8C6ECDF145FDCC697D6B72562FA2E0418B2B80A038D0DC3B769EB027484087CCE6652488D2B3816742AC9C2355B17411C47EACDD7467566B302F512806E331FAD964BF000169F641DF040103DF0300BF010131DF070101", "VF_9F0605A0000000659F22010FDF050420291231DF0281909EFBADDE4071D4EF98C969EB32AF854864602E515D6501FDE576B310964A4F7C2CE842ABEFAFC5DC9E26A619BCF2614FE07375B9249BEFA09CFEE70232E75FFD647571280C76FFCA87511AD255B98A6B577591AF01D003BD6BF7E1FCE4DFD20D0D0297ED5ECA25DE261F37EFE9E175FB5F12D2503D8CFB060A63138511FE0E125CF3A643AFD7D66DCF9682BD246DDEA1DF040103DF03142A1B82DE00F5F0C401760ADF528228D3EDE0F403BF010131DF070101", "VF_9F0605A0000000659F220113DF050420250101DF0281F8A3270868367E6E29349FC2743EE545AC53BD3029782488997650108524FD051E3B6EACA6A9A6C1441D28889A5F46413C8F62F3645AAEB30A1521EEF41FD4F3445BFA1AB29F9AC1A74D9A16B93293296CB09162B149BAC22F88AD8F322D684D6B49A12413FC1B6AC70EDEDB18EC1585519A89B50B3D03E14063C2CA58B7C2BA7FB22799A33BCDE6AFCBEB4A7D64911D08D18C47F9BD14A9FAD8805A15DE5A38945A97919B7AB88EFA11A88C0CD92C6EE7DC352AB0746ABF13585913C8A4E04464B77909C6BD94341A8976C4769EA6C0D30A60F4EE8FA19E767B170DF4FA80312DBA61DB645D5D1560873E2674E1F620083F30180BD96CA589DF040103DF031454CFAE617150DFA09D3F901C9123524523EBEDF3BF010131DF070101", "VF_9F0605A0000001529F22015CDF050420291231DF0281B0833F275FCF5CA4CB6F1BF880E54DCFEB721A316692CAFEB28B698CAECAFA2B2D2AD8517B1EFB59DDEFC39F9C3B33DDEE40E7A63C03E90A4DD261BC0F28B42EA6E7A1F307178E2D63FA1649155C3A5F926B4C7D7C258BCA98EF90C7F4117C205E8E32C45D10E3D494059D2F2933891B979CE4A831B301B0550CDAE9B67064B31D8B481B85A5B046BE8FFA7BDB58DC0D7032525297F26FF619AF7F15BCEC0C92BCDCBC4FB207D115AA65CD04C1CF982191DF040103DF0314C165C48EB36DDF969DDC0B326312AFE2F6B52713BF010131DF070101", "VF_9F0605A0000001529F22015BDF050420291231DF028190D3F45D065D4D900F68B2129AFA38F549AB9AE4619E5545814E468F382049A0B9776620DA60D62537F0705A2C926DBEAD4CA7CB43F0F0DD809584E9F7EFBDA3778747BC9E25C5606526FAB5E491646D4DD28278691C25956C8FED5E452F2442E25EDC6B0C1AA4B2E9EC4AD9B25A1B836295B823EDDC5EB6E1E0A3F41B28DB8C3B7E3E9B5979CD7E079EF024095A1D19DDDF040103DF03140000000000000000000000000000000000000000BF010131DF070101", "VF_9F0605A0000001529F22015DDF050420241231DF0281F8AD938EA9888E5155F8CD272749172B3A8C504C17460EFA0BED7CBC5FD32C4A80FD810312281B5A35562800CDC325358A9639C501A537B7AE43DF263E6D232B811ACDB6DDE979D55D6C911173483993A423A0A5B1E1A70237885A241B8EEBB5571E2D32B41F9CC5514DF83F0D69270E109AF1422F985A52CCE04F3DF269B795155A68AD2D6B660DDCD759F0A5DA7B64104D22C2771ECE7A5FFD40C774E441379D1132FAF04CDF55B9504C6DCE9F61776D81C7C45F19B9EFB3749AC7D486A5AD2E781FA9D082FB2677665B99FA5F1553135A1FD2A2A9FBF625CA84A7D736521431178F13100A2516F9A43CE095B032B886C7A6AB126E203BE7DF040103DF0314B51EC5F7DE9BB6D8BCE8FB5F69BA57A04221F39BBF010131DF070101", "VF_9F0605A000000003 9F220194 DF0403000003 DF0314C4A3C43CCF87327D136B804160E47D43B60E6E0F DF050420491231 DF0281F8ACD2B12302EE644F3F835ABD1FC7A6F62CCE48FFEC622AA8EF062BEF6FB8BA8BC68BBF6AB5870EED579BC3973E121303D34841A796D6DCBC41DBF9E52C4609795C0CCF7EE86FA1D5CB041071ED2C51D2202F63F1156C58A92D38BC60BDF424E1776E2BC9648078A03B36FB554375FC53D57C73F5160EA59F3AFC5398EC7B67758D65C9BFF7828B6B82D4BE124A416AB7301914311EA462C19F771F31B3B57336000DFF732D3B83DE07052D730354D297BEC72871DCCF0E193F171ABA27EE464C6A97690943D59BDABB2A27EB71CEEBDAFA1176046478FD62FEC452D5CA393296530AA3F41927ADFE434A2DF2AE3054F8840657A26E0FC617 BF010131 DF070101", "VF_9F0605A000000333 9F22010A DF050420301231 DF0280B2AB1B6E9AC55A75ADFD5BBC34490E53C4C3381F34E60E7FAC21CC2B26DD34462B64A6FAE2495ED1DD383B8138BEA100FF9B7A111817E7B9869A9742B19E5C9DAC56F8B8827F11B05A08ECCF9E8D5E85B0F7CFA644EFF3E9B796688F38E006DEB21E101C01028903A06023AC5AAB8635F8E307A53AC742BDCE6A283F585F48EF DF0314C88BE6B2417C4F941C9371EA35A377158767E4E3 DF040103 BF010131 DF070101" ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
-- CtlsApplicationConfigShared
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'ctlsApplicationConfigShared' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'multiDeviceEmv')), '[ "XAC_DF1006999999999999 DF1106000000000000 DF1206000000050000 DF130100 DF190113 DF220100", "VF_DF1906000000000000 DF2006999999999999 DF2106000000050000" ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
-- contactApplicationConfigShared
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'contactApplicationConfigShared' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'multiDeviceEmv')), '[ "XAC_5F2A020784 5F360102 9F0106012345678901 9F15021234 9F160F313233343536373839303132333435 9F1B0400000000 9F1D080000000000000000 9F350122 9F370400000000 9F3C020784 9F3D0102 DF100101 DF110400000002 DF12069F02069F0306 DF13159F02069F03069F1A0295055F2A029A039C019F3704 DF170100 DF180100 DF190400000000", "VF_9F1B0400000000 9F350122 DF14039F3704 DF150400000000 DF160101 DF170101 DF010100" ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
-- ctlsApplicationConfigs
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'ctlsApplicationConfigs' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'multiDeviceEmv')), '[ "XAC_DF0B0100 C00102 9F0607A0000000031010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000", "XAC_DF0B0101 C00102 9F0607A0000000032010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000", "XAC_DF0B0102 C00102 9F0607A0000000033010 DF0C0141 DF0D050010000000 DF0E05DC4004F800 DF0F05DC4000A800 500456534443 9F09020105 9F15020011 9F1D084400800000000000", "XAC_DF0B0103 C00101 9F0607A0000000041010 DF0C0109 DF0D050000000000 DF0E05FC50BCF800 DF0F05FC50BCA000 500A4D617374657243617264 9F09020002 9F15020011 9F1D084400800000000000 E60F E7039C0120 E808DF0D05FFFFFFFFFF DF2606999999999999", "XAC_DF0B0104 C00101 9F0607A0000000043060 DF0C010B DF0D050000800000 DF0E05FC50BCF800 DF0F05FC50BCA000 50074D61657374726F 9F09020002 9F15020011 9F1D084400800000000000 E60F E7039C0120 E808DF0D05FFFFFFFFFF DF2606999999999999", "XAC_DF0B0105 C0010C 9F0607A0000000651010 DF0C0143 DF0D050050000000 DF0E05FC60ACF800 DF0F05FC6024A800 50034A4342 9F09020200 9F15020011 9F1D084400800000000000", "VF_9F0607A0000000031010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000", "VF_9F0607A0000000032010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000", "VF_9F0607A0000000033010 9F09020105 DF1105DC4000A800 DF1205DC4004F800 DF13050010000000", "VF_9F0607A0000000041010 9F09020002 DF1105FC50BCA000 DF1205FC50BCA000 DF13050000000000 DF811B0190", "VF_9F0607A0000000043060 9F09020002 DF1105FC50BCA000 DF1205FC50BCA000 DF13050000800000 DF811B0190" ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
-- contactApplicationConfigs
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'contactApplicationConfigs' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'multiDeviceEmv')), '[ "XAC_9F0607A0000000031010 9F0902008C DF1405DC4000A800 DF15050050000000 DF1605DC4004F800", "XAC_9F0607A0000000032010 9F0902008C DF1405DC4000A800 DF15050050000000 DF1605DC4004F800", "XAC_9F0607A0000000033010 9F0902008C DF1405DC4000A800 DF15050050000000 DF1605DC4004F800", "XAC_9F0607A0000000041010 9F09020002 DF1405FC50BCA000 DF15050050000000 DF1605FC50BCF800 E60F E7039C0120 E808DF0D05FFFFFFFFFF", "XAC_9F0607A0000000042203 9F09020002 DF1405FC50BCA000 DF15050050000000 DF1605FC50BCF800 E60F E7039C0120 E808DF0D05FFFFFFFFFF", "XAC_9F0607A0000000043060 9F09020002 DF1405FC50BCA000 DF15050050000000 DF1605FC50BCF800 E60F E7039C0120 E808DF0D05FFFFFFFFFF", "XAC_9F0607A0000003330101 9F09020020 DF1405D84000A800 DF15050050000000 DF1605D84004F800", "XAC_9F0606A00000002501 9F09020001 DF1405CC00FC8000 DF15050010000000 DF1605DE00FC9800", "XAC_9F0607A0000000651010 9F09020200 DF1405FC6024A800 DF15050050000000 DF1605FC60ACF800", "XAC_9F0607A0000001523010 9F09020001 DF1405DC00002000 DF15050010000000 DF1605FCE09CF800", "XAC_9F0607A0000001524010 9F09020001 DF1405DC00002000 DF15050010000000 DF1605FCE09CF800", "XAC_9F0607A0000005291010 9F09020001 DF1405FFFFFFFFFF DF15050010000000 DF1605FFFFFFFFFF", "VF_9F0607A0000000031010 9F0902008C DF1105DC4000A800 DF1205DC4004F800 DF13050050000000", "VF_9F0607A0000000032010 9F0902008C DF1105DC4000A800 DF1205DC4004F800 DF13050050000000", "VF_9F0607A0000000033010 9F0902008C DF1105DC4000A800 DF1205DC4004F800 DF13050050000000", "VF_9F0607A0000000041010 9F09020002 DF1105FC50BCA000 DF1205FC50BCF800 DF13050050000000", "VF_9F0607A0000000042203 9F09020002 DF1105FC50BCA000 DF1205FC50BCF800 DF13050050000000", "VF_9F0607A0000000043060 9F09020002 DF1105FC50BCA000 DF1205FC50BCF800 DF13050050000000" ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);
-- update card definitions
INSERT into profile_data (profile_id, data_element_id, datavalue, version, updated_at, updated_by, created_at, created_by, approved, overriden, is_encrypted) values (@profile_id, (SELECT data_element_id FROM data_element WHERE name = 'cardDefinitions' AND data_group_id = (SELECT data_group_id FROM data_group WHERE name = 'emv')), '[ { "cardName":"MASTER", "brandCode":1, "minLength":12, "maxLength":19, "tacDenial":"0050000000", "tacOnline":"FC50BCF800", "ranges":[ "2221>2720", "51>55" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "dccCtlsLimitCheckPP": 500.00, "dccCtlsLimitCheckFX": 500.00, "dccEnabled":true, "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>", "dccPrintFee":false }, { "cardName":"VISA", "brandCode":2, "minLength":16, "maxLength":19, "tacDenial":"0050000000", "tacOnline":"DC4004F800", "ranges":[ "4" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "dccCtlsLimitCheckPP": 0, "dccCtlsLimitCheckFX": 0, "dccEnabled":true, "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>", "dccPrintFee":true }, { "cardName":"AMEX", "brandCode":3, "amexMid":true, "minLength":1, "maxLength":19, "tacDenial":"0010000000", "tacOnline":"DE00FC9800", "ranges":[ "34", "37" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "disableOnlinePin":true }, { "cardName":"DINERS", "brandCode":4, "minLength":14, "maxLength":19, "tacDenial":"0010000000", "tacOnline":"FCE09CF800", "ranges":[ "36", "6510" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E 9F71", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E 9F71" }, { "cardName":"UNIONPAY", "brandCode":11, "minLength":16, "maxLength":19, "refundsEnabled":false, "tacDenial":"0050000000", "tacOnline":"D84004F800", "ranges":[ "601382", "601428", "602907", "602969", "603265", "603367", "603601", "603694", "603708", "620000>621799", "621977", "622126>626999", "626202", "627066>627067", "628200>628899", "629100>629399", "6858", "69075", "90>98", "816399", "817199", "990027", "998800>998802" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E" }, { "cardName":"JCB", "brandCode":12, "minLength":16, "maxLength":16, "tacDenial":"0050000000", "tacOnline":"FC60ACF800", "ranges":[ "3528>3589" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E" }, { "cardName":"MERCURY", "brandCode":0, "minLength":16, "maxLength":19, "tacDenial":"0010000000", "tacOnline":"FFFFFFFFFF", "ranges":[ "650401", "650008", "650483", "656001", "978432>978439" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E" }, { "cardName":"MAESTRO", "brandCode":1, "minLength":12, "maxLength":19, "tacDenial":"0050000000", "tacOnline":"FC50BCF800", "ranges":[ "500000>509999", "560000>699999" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "dccEnabled":true, "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>", "dccPrintFee":false }, { "cardName":"NSPK", "brandCode":14, "minLength":16, "maxLength":19, "tacDenial":"0010000000", "tacOnline":"B6408C8000", "ranges":[ "2201" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E 9F4C 9F40", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E 9F4C 9F40 " }, { "cardName":"RUPAY", "brandCode":1, "minLength":16, "maxLength":19, "tacDenial":"0000000000", "tacOnline":"FEF8FCF8F0", "ranges":[ "508500>508999", "606985>607984", "608001>608500", "817200>820199" ], "offlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F07 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "onlineChipData":"50 57 82 84 89 95 5A 5F24 5F25 5F2A 5F34 8A 9A 9B 9C 9F01 9F02 9F03 9F06 9F07 9F09 9F10 9F11 9F12 9F15 9F16 9F1A 9F1E 9F21 9F26 9F27 9F33 9F34 9F35 9F36 9F37 9F39 9F40 9F41 9F53 9F6E", "dccEnabled":false, "dccDisclaimer":"I ACCEPT THAT I HAVE BEEN GIVEN A CHOICE OF CURRENCIES FOR PAYMENT AND THAT THIS CHOICE IS FINAL. I ACCEPT THE CONVERSION RATE, THE FINAL AMOUNT AND THE SELECTED TRANSACTION. THIS CURRENCY CONVERSION SERVICE IS PROVIDED BY <<MERCHANT NAME>>", "dccPrintFee":false } ]', 1, NOW(), 'NISuper', NOW(), 'NISuper', 1, 1, 0) ON DUPLICATE KEY UPDATE datavalue=VALUES(datavalue), updated_at = VALUES(updated_at);