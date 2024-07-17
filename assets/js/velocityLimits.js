rivets.binders.i = (el, val) => {
    const jEl = $(el);

    let part = jEl.attr('rv-value');
    if (part == null) {
        part = jEl.attr('rv-on-click')
    }

    jEl.attr('id', val + "-" + part.split('.').pop())
};

$(document).ready(function () {
    getVelocityLimits();
    numbersOnly()
});
let limits;
let txnLimitError = "The chosen limit type for this transaction type is already set";
let txnUpdateSuccess = "Transaction limit updated";
let schemeLimitError = "Duplicate scheme limits detected";
let schemeNameError = "No scheme name set";
let limitSaveSuccess = "Velocity limits saved";
let buttonId;

let limitData = {
    LimitIndex: '',
    TxnLimits: []
};

function getVelocityLimits() {
    $.ajax({
        data: {
            tidId: -1,
            siteId: $("#siteID")[0].value,
            profileId: $("#profileID")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: 'POST',
        url: '/GetSiteVelocityLimits',

        success: function (data) {
            if (data != null) {
                limits = data;
                //Add the start modal function for site level
                limits.SiteLimits.showSiteTxnLimitModal = startSiteTxnLimitModal;
                // When the rivet is bound this will set the id of scheme agnostic limits to site-BatchLimit etc
                limits.SiteLimits.prefix = "site";

                tidyNonLimitedLimits(limits.SiteLimits);
                tidySchemeLimits();

                rivets.bind($('#siteVelocityLimitBody'), { SiteLimits: limits.SiteLimits });
                rivets.bind($('#velocityLimitBody'), { limits: limits.Limits });
                rivets.bind($('#txnSelectField'), { availableTransactions: limits.AvailableTransactions });
                rivets.bind($('#limitSelectField'), { availableLimits: limits.AvailableLimits });

                displayOverridenLimits(limits, "site");
            }
        },
        error: function (data) {
        }
    })
}

function tidySchemeLimits() {
    //Add the limit function for scheme level
    if (limits.Limits == null) {
        limits.Limits = [];
    } else {
        limits.Limits.forEach(function (limit) {
            addLimitMethods(limit);
            tidyNonLimitedLimits(limit);
            // When the rivet is bound this will set the id of scheme limits to 0-site-BatchLimit etc
            limit.prefix = limit.Index + "-site";
        });

        // Sort the limits by their index
        limits.Limits.sort((a, b) => a.Index - b.Index);
    }
}

function displayOverridenLimits(limits, identifier) {
    var site = limits.SiteLimits;
    var schemes = limits.Limits;
    //Check to see if there are site level transaction limits set and if so set to green
    if (site.TxnLimits != null) {
        if (site.TxnLimits.length > 0) {
            $('#' + identifier + '-showSiteTxnLimitModal').addClass('activated-btn');
        }
    }

    //Do the same for scheme level transactions
    for (x of schemes) {
        if (x.TxnLimits != null) {
            if (x.TxnLimits.length > 0) {
                let cleanIdentifier = DOMPurify.sanitize(identifier);
                let cleanIndex = DOMPurify.sanitize(x.Index);
                $('#' + Number(cleanIndex) + '-' + cleanIdentifier + '-showTxnLimitModal').addClass('activated-btn');
            }
        }
    }
}

function tidyNonLimitedLimits(limit) {
    if (limit.DailyCount === -1) {
        limit.DailyCount = ""
    }
    if (limit.BatchCount === -1) {
        limit.BatchCount = ""
    }
    if (limit.SingleTransLimit === -1) {
        limit.SingleTransLimit = ""
    }
    if (limit.DailyLimit === -1) {
        limit.DailyLimit = ""
    }
    if (limit.BatchLimit === -1) {
        limit.BatchLimit = ""
    }
}

function addRows() {
    hideWarning();
    let numberOfRows = limits.Limits.length
    let limit = {
        ID: newGuid(),
        Scheme: '',
        DailyCount: null,
        DailyLimit: null,
        BatchCount: null,
        BatchLimit: null,
        SingleTransLimit: null,
        TxnLimits: [],
        Index: numberOfRows,
        removeLimit: deleteRow,
        prefix: numberOfRows + "-site"
    };
    addLimitMethods(limit);
    limits.Limits.push(limit);
    numbersOnly()
}

function closeTxnModal() {
    $("#limitModal").modal("hide");
    saveTxnLimits(limitData.LimitIndex, false);
    if (limitData.TxnLimits.length > 0) {
        $('#' + buttonId).addClass('activated-btn')
    }
}

function deleteRow(event, limit) {
    hideWarning();
    limits.Limits.splice(limit.index, 1);
}

function saveTxnLimits(index, tidLimit) {

    if (index === '') {
        return
    } else if (index === 'site') {
        if (tidLimit) {
            tidLimits.SiteLimits.TxnLimits = tidLimitData.TxnLimits
        } else {
            limits.SiteLimits.TxnLimits = limitData.TxnLimits
        }
    }

    if ((typeof limits?.Limits[index] !== 'undefined') || (typeof tidLimits?.Limits[index] !== 'undefined')) {

        if (tidLimit) {
            tidLimits.Limits[index].TxnLimits = tidLimitData.TxnLimits
        } else {
            limits.Limits[index].TxnLimits = limitData.TxnLimits
        }
    }
}

function startSiteTxnLimitModal(event, limit) {

    limitData.LimitIndex = "site";
    limitData.TxnLimits = limits.SiteLimits.TxnLimits || [];
    buttonId = this.id;
    displayTxnLimits();
    $('#limitModal').modal('show')
}

function startTxnLimitModal(event, bindings) {
    let limitClicked = bindings.limit;
    // save current filter data
    saveTxnLimits(limitData.LimitIndex, false);
    // load alt data
    limitData.LimitIndex = limitClicked.Index;

    limitClicked.Index = bindings.index;

    limitData.TxnLimits = limits.Limits[limitClicked.Index].TxnLimits || [];

    buttonId = this.id;
    displayTxnLimits();
    $('#limitModal').modal('show');
}

function saveSiteLimits() {
    hideWarning();
    let tidId = -1;
    let limitJson = [];

    saveTxnLimits(limitData.LimitIndex, false);
    limitJson.push(buildSiteSave(limits.SiteLimits, 3));
    let limitModel = buildSchemeSave(limits.Limits, tidId, 1);
    let flagStatus = $("#hidden_velocityFlagStatusOption").val();
    if (limitModel.error === false) {
        let schemeLimits = limitModel.limits
        for (x of schemeLimits) {
            limitJson.push(x)
        }

        let dataObj = new FormData()
        if (flagStatus=="file"){
            var uploadFile = $('#flagged-tids-file')[0].files[0];
            dataObj.append("flagged-tids-file",  uploadFile);
        }
        dataObj.append("tidId", tidId);
        dataObj.append("siteId", $("#siteID")[0].value);
        dataObj.append("profileId", $("#profileID")[0].value);
        dataObj.append("limitLevel", 3);
        dataObj.append("csrfmiddlewaretoken", $("input[name=csrfmiddlewaretoken]")[0].value);
        dataObj.append("limits",limitJson);
        dataObj.append("flagStatus", flagStatus);
        ////TODO: Clean this up to be more generic and get key value pairs from the table
        dataObj.append("dailyTxnCleanseTime", $("#site-fraud-group-table #core-dailyTxnCleanseTime").val());
        $("#flagging-checkbox").find("input:checkbox:checked").not("[disabled]").each(function() {
            dataObj.append($(this).attr('id'),$(this).val());
        });
        $.ajax({
            type: 'POST',
            url: '/SaveSiteVelocityLimits',
            data: dataObj,
            processData: false,
            contentType: false,
            success: function (data) {
                displayWarningMessage(limitSaveSuccess)
            },
            error: function (data) {
                if (data.getResponseHeader("content-type") === "application/json") {
                    let validationMessages = JSON.parse(data.responseText).join('<br/>');
                    displayWarningMessage(validationMessages)
                } else {
                    displayWarningMessage(data.responseText)
                }
            }
        });
    }

    //Ensure site level limits that have no restriction display as blank
    tidyNonLimitedLimits(limits.SiteLimits)
    if (limits.Limits == null) {
        limits.Limits = []
    } else {
        limits.Limits.forEach(function (limit) {
            addLimitMethods(limit);
            tidyNonLimitedLimits(limit);
        });
    }
}

function buildSiteSave(siteLimit, limitLevel) {
    //For new sites the site level limit ID will be blank, here we can set it
    if (siteLimit.ID === "") {
        siteLimit.ID = newGuid()
    }

    setNullToNegativeOne(siteLimit);
    delete siteLimit.showTxnLimitModal;

    return JSON.stringify({
        ID: siteLimit.ID,
        Scheme: siteLimit.Scheme,
        DailyCount: parseInt(siteLimit.DailyCount, 10),
        DailyLimit: parseInt(siteLimit.DailyLimit, 10),
        BatchCount: parseInt(siteLimit.BatchCount, 10),
        BatchLimit: parseInt(siteLimit.BatchLimit, 10),
        SingleTransLimit: parseInt(siteLimit.SingleTransLimit, 10),
        TxnLimits: siteLimit.TxnLimits,
        Level: limitLevel
    });
}

function buildSchemeSave(schemeLimits, tidId, limitLevel) {
    let limitModel = [];
    if (validateLimits(schemeLimits, tidId)) {
        let limitJson = [];
        schemeLimits.forEach(function (limit) {
            removeMethods(limit);
            setNullToNegativeOne(limit)
            limitJson.push(JSON.stringify({
                ID: limit.ID,
                Scheme: limit.Scheme,
                DailyCount: parseInt(limit.DailyCount, 10),
                DailyLimit: parseInt(limit.DailyLimit, 10),
                BatchCount: parseInt(limit.BatchCount, 10),
                BatchLimit: parseInt(limit.BatchLimit, 10),
                SingleTransLimit: parseInt(limit.SingleTransLimit, 10),
                TxnLimits: limit.TxnLimits,
                Level: limitLevel,
                Index: limit.Index
            }));
            tidyNonLimitedLimits(limit);
        });
        limitModel.limits = limitJson
        limitModel.error = false
        return limitModel
    } else {
        limitModel.error = true
        return limitModel
    }
}

function validateLimits(schemeLimits, tidId) {

    var unsavedLimits = [];

    for (x of schemeLimits) {
        if (unsavedLimits.includes(x.Scheme)) {
            if (tidId === -1) {
                displayWarningMessage(schemeLimitError)
            } else {
                displayWarningMessage(schemeLimitError, "tidFraudModalBody")
            }
            return false
        } else if (x.Scheme === "") {
            if (tidId === -1) {
                displayWarningMessage(schemeNameError);
            } else {
                displayWarningMessage(schemeNameError, "tidFraudModalBody");
            }
            return false
        }
        unsavedLimits.push(x.Scheme);
        setNullToNegativeOne(x)
    }

    return true
}

function setNullToNegativeOne(limit) {
    if (limit.DailyCount === "" || limit.DailyCount === null) {
        limit.DailyCount = -1
    }
    if (limit.BatchCount === "" || limit.BatchCount === null) {
        limit.BatchCount = -1
    }
    if (limit.SingleTransLimit === "" || limit.SingleTransLimit === null) {
        limit.SingleTransLimit = -1
    }
    if (limit.DailyLimit === "" || limit.DailyLimit === null) {
        limit.DailyLimit = -1
    }
    if (limit.BatchLimit === "" || limit.BatchLimit === null) {
        limit.BatchLimit = -1
    }
}

function addLimitMethods(limit) {
    limit.removeLimit = deleteRow;
    limit.showTxnLimitModal = startTxnLimitModal
}

function removeMethods(limit) {
    delete limit.removeLimit;
    delete limit.showTxnLimitModal
}

function addTxnLimit() {
    let txnType = $('#txnSelectField').find(':selected');
    let txnReadable = $('#txnSelectField option:selected').text();
    let limitType = $('#limitSelectField').find(':selected');
    let value = parseInt($('#limitValueField').val());
    let override = $('#override').val();

    if (validateTxnLimit(txnType.val(), limitType.val(), value, override)) {

        if (override === 'true') {
            tidLimitData.TxnLimits.push({
                TxnLimitID: newGuid(),
                TxnType: txnType.val(),
                TxnTypeReadable: txnReadable,
                LimitType: limitType.val(),
                Value: value
            });
            displayTIDTxnLimits()
        } else {
            limitData.TxnLimits.push({
                TxnLimitID: newGuid(),
                TxnType: txnType.val(),
                TxnTypeReadable: txnReadable,
                LimitType: limitType.val(),
                Value: value
            });
            displayTxnLimits()
        }
    }
}

//Check that the chosen combination of txn type and limit type is not already set for this scheme
function validateTxnLimit(txnType, limitType, value, override) {

    var limits
    if (override === "true") {
        limits = tidLimitData.TxnLimits
    } else {
        limits = limitData.TxnLimits
    }

    for (x of limits) {
        if (x.TxnType === txnType && x.LimitType === limitType) {
            displayWarningMessage(txnLimitError, "txnLimitModal");
            return false
        }
    }

    return true
}

function displayTxnLimits() {
    $.ajax({
        data: {
            override: false,
            siteId: $("#siteID")[0].value,
            profileId: $("#profileID")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
            TxnLimitData: JSON.stringify(limitData)
        },
        type: 'POST',
        url: '/velocityLimitRow',
        success: function (data) {
            $('#txnLimitContainer').html(data);
            disableUpdateButtons();
            numbersOnly()
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function setTxnLimitValue(limitId, value, override) {
    let txnLimit = findTxnLimit(tidOrSite(override), limitId);
    txnLimit.Value = value
}

function setLimitType(limitId, limitType, override) {
    let txnLimit = findTxnLimit(tidOrSite(override), limitId);
    txnLimit.LimitType = limitType
}

//determines whether the data to be passed is tid or site level
function tidOrSite(override) {

    if (override === "true") {
        return tidLimitData
    } else {
        return limitData
    }
}

function findTxnLimit(group, limitId) {
    let target = null;
    group.TxnLimits.forEach(function (element) {
        if (element.TxnLimitID === limitId) {
            target = element
        }
    });
    return target
}

function removeTxnLimit(id) {
    let override = $('#override').val();

    if (override === "true") {
        removeTxnLimitFromGroup(tidLimitData, id)
    } else {
        removeTxnLimitFromGroup(limitData, id)
    }
}

function removeTxnLimitFromGroup(group, id) {
    group.TxnLimits = group.TxnLimits.filter(function (element) {
        return element.TxnLimitID !== id
    })
}

function newGuid() {
    return (S4() + S4() + '-' + S4() + '-4' + S4().substr(0, 3) + '-' + S4() + '-' + S4() + S4() + S4()).toLowerCase()
}

function S4() {
    return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1)
}

function updateTxnLimits(txnLimitID, rowID) {
    let txnType = document.getElementById('txnType-' + rowID).textContent;
    // NEX-10954 For some of the transaction type which has space 
    //eg. Tip Sale added div element  
    let limitType = document.getElementById('limitType-' + rowID).value;
    let value = parseInt(document.getElementById('value-' + rowID).value, 10);
    let override = $('#override').val();
    if (validateTxnLimit(txnType, limitType, value, override)) {
        setTxnLimitValue(txnLimitID, value, override);
        setLimitType(txnLimitID, limitType, override);
        displayWarningMessage(txnUpdateSuccess, "txnLimitModal");
    }
}

//Enables txn update button when values are amended
function showUpdateButton(txnLimitID) {
    $('#update-txnLimit-' + txnLimitID).prop('disabled', false)
}

//Disables all txn limit update buttons until their values are amended
function disableUpdateButtons() {
    for (x of limitData.TxnLimits) {
        const limitID = x.TxnLimitID;
        $('#update-txnLimit-' + limitID).prop('disabled', true)
    }
}

//Prevents the use of non numeric characters in number fields
function numbersOnly() {
    $('input.numbersOnly').on("keypress", function (e) {

        var keyValue = String.fromCharCode(e.keyCode);
        //prevent: "e", "=", ",", "-", "."
        if (keyValue === "e" || keyValue === "=" || keyValue === "," || keyValue === "-" || keyValue === ".") {
            e.preventDefault();
        }

        var limit = this.value;
        if (limit >= 99999999) {
            e.preventDefault()
        }
    });
}


