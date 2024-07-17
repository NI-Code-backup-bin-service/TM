
let tidLimits
let tidLimitData = {
    LimitIndex: '',
    TxnLimits: []
};
var tidSiteVelocityLimitBody
var tidVelocityLimitBody
var txnSelectField
var limitSelectField
let tidButtonId;

function getTidVelocityLimits(valuesSaved) {
    $.ajax({
        data: {
            tidId: $("#tidFraudID")[0].value,
            siteId: $("#siteID")[0].value,
            profileId: $("#profileID")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: 'POST',
        url: '/GetSiteVelocityLimits',

        success: function (data) {
            if (data != null) {
                bindTidRivets();
                tidLimits = data;
                let tidID = $("#tidFraudID")[0].value;

                //Add the start modal function for site level
                tidLimits.SiteLimits.showSiteTxnLimitModal = startTIDSiteTxnLimitModal;
                // When the rivet is bound this will set the id of scheme agnostic limits to 80800808-BatchLimit etc
                tidLimits.SiteLimits.prefix = tidID;

                tidyNonLimitedLimits(tidLimits.SiteLimits);
                tidyTidSchemeLimits(tidID);
                updateRivetFraudElements(tidID, valuesSaved);
                displayOverridenLimits(tidLimits, tidID);
            }
        },
        error: function (data) {
        }
    })
}

function bindTidRivets() {
    rivets.binders.tid = (el, val) => {
        const jEl = $(el);
        let part = jEl.attr('rv-value');
        if (part == null) {
            part = jEl.attr('rv-on-click');
        }

        jEl.attr('id', val + "-" + part.split('.').pop());
    };
}

function updateRivetFraudElements(tidID, valuesSaved) {
    if (valuesSaved === true) {
        tidSiteVelocityLimitBody = rivets.bind($('#' + tidID + '-tidSiteVelocityLimitBody'), { SiteLimits: tidLimits.SiteLimits });
        tidVelocityLimitBody.update(tidLimits.Limits);
        txnSelectField.update(tidLimits.AvailableTransactions);
        limitSelectField.update(tidLimits.AvailableLimits);
    } else {
        tidSiteVelocityLimitBody = rivets.bind($('#' + tidID + '-tidSiteVelocityLimitBody'), { SiteLimits: tidLimits.SiteLimits });
        tidVelocityLimitBody = rivets.bind($('#' + tidID + '-tidVelocityLimitBody'), { limits: tidLimits.Limits });
        txnSelectField = rivets.bind($('#txnSelectField'), { availableTransactions: tidLimits.AvailableTransactions });
        limitSelectField = rivets.bind($('#limitSelectField'), { availableLimits: tidLimits.AvailableLimits });
    }
}

function tidyTidSchemeLimits(tidID) {
    //Add the limit function for scheme level
    if (tidLimits.Limits == null) {
        tidLimits.Limits = [];
    } else {
        tidLimits.Limits.forEach(function (limit) {
            addTIDLimitMethods(limit);
            tidyNonLimitedLimits(limit);

            // When the rivet is bound this will set the id of scheme limits to 0-80800808-BatchLimit etc
            limit.prefix = limit.Index + "-" + tidID;
        });

        // Sort the limits by their index
        tidLimits.Limits.sort((a, b) => a.Index - b.Index);
    }
}

function saveTidSiteLimits() {
    hideWarning("tidFraudModalBody");
    let tidId = $("#tidFraudID")[0].value;
    let limitJson = [];

    saveTxnLimits(tidLimitData.LimitIndex, true);
    limitJson.push(buildSiteSave(tidLimits.SiteLimits, 4));
    let limitModel = buildSchemeSave(tidLimits.Limits, tidId, 2);

    if (limitModel.error === false) {
        let schemeLimits = limitModel.limits
        for (x of schemeLimits) {
            limitJson.push(x)
        }

        $.ajax({
            type: 'POST',
            url: '/SaveSiteVelocityLimits',
            data: {
                tidId: tidId,
                siteId: $("#siteID")[0].value,
                profileId: $("#profileID")[0].value,
                limitLevel: 4,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
                limits: limitJson.toString(),
                //TODO: Clean this up to be more generic and get key value pairs from the table
                dailyTxnCleanseTime: $("#tid-fraud-group-table #core-dailyTxnCleanseTime").val()
            },
            success: function (data) {
                displayWarningMessage(limitSaveSuccess, "tidFraudModalBody");

            },
            error: function (data) {
                //the data is a json array of messages
                if (data.getResponseHeader("content-type") === "application/json") {
                    let validationMessages = JSON.parse(data.responseText).join('<br/>');
                    displayWarningMessage(validationMessages, "tidFraudModalBody")

                } else {
                    displayWarningMessage(data.responseText, "tidFraudModalBody")

                }
            }
        });
    }
    tidyNonLimitedLimits(tidLimits.SiteLimits);
    tidyTidSchemeLimits(tidId);
    updateRivetFraudElements(tidId, true);
}

function deleteTidSiteLimits() {
    tidLimits.SiteLimits = Object.assign(tidLimits.SiteLimits, { DailyCount: "", DailyLimit: "", BatchCount: "", BatchLimit: "", SingleTransLimit: "" });
}

function deleteTidFraudLimits() {

    confirmDialog("Delete Override?", 'Are you sure you want to delete the Fraud Override?', function () {

        $.ajax({
            url: "/deleteSiteVelocityLimits",
            method: "POST",
            data: {
                tidId: $("#tidFraudID")[0].value,
                siteId: $("#siteID")[0].value,
                profileId: $("#profileID")[0].value,
                limitLevel: 4,
                csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
            },
            success: function (d) {
                closeTIDFraudModal()
            },
            error: function (data) {
                displayWarningMessage(data.responseText, "tidFraudModalBody")
            }
        });
    })
}

function startTIDSiteTxnLimitModal(event, limit) {

    tidLimitData.LimitIndex = "site";
    tidLimitData.TxnLimits = tidLimits.SiteLimits.TxnLimits || [];
    tidButtonId = this.id;
    displayTIDTxnLimits();
    $('#TIDlimitModal').modal('show')
}

function startTIDTxnLimitModal(event, bindings) {

    let limitClicked = bindings.limit;

    // save current filter data
    saveTxnLimits(tidLimitData.LimitIndex, true);
    // load alt data
    tidLimitData.LimitIndex = limitClicked.Index;

    limitClicked.Index = bindings.index;

    tidLimitData.TxnLimits = tidLimits.Limits[limitClicked.Index].TxnLimits || [];

    tidButtonId = this.id;
    displayTIDTxnLimits();
    $('#TIDlimitModal').modal('show')
}

function addTIDLimitMethods(limit) {
    limit.removeLimit = deleteTidRow;
    limit.showTxnLimitModal = startTIDTxnLimitModal
}

function deleteTidRow(event, limit) {
    hideWarning();
    tidLimits.Limits.splice(limit.index, 1);
}

function addTIDSchemeRows(tidID) {
    hideWarning();
    let numberOfRows = tidLimits.Limits.length;
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
        prefix: numberOfRows + "-" + tidID
    };
    addTIDLimitMethods(limit);
    tidLimits.Limits.push(limit);
    numbersOnly();
}

function closeTIDTxnModal() {

    $("#TIDlimitModal").modal("hide");
    saveTxnLimits(tidLimitData.LimitIndex, true);
    if (tidLimitData.TxnLimits.length > 0) {
        for (x of tidLimitData.TxnLimits) {
            x.TxnLimitID = newGuid()
        }
        $(tidButtonId).addClass('activated-btn')
    }
}

function closeTIDFraudModal() {
    tidSiteVelocityLimitBody.unbind();
    tidVelocityLimitBody.unbind();
    txnSelectField.unbind();
    limitSelectField.unbind();
    $("#tidFraudModal").modal("hide");
    $.ajax({
        data: {
            profileId: $("#profileID")[0].value,
            pageSize: $("#pageSize")[0].value,
            pageNumber: $("#pageNumber")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        type: 'POST',
        url: '/TidFraudClose',
        success: function (data) {
            $('#tids').html(sanitizeHTML(data));
            hideWarning();
            $('[name^=multi]').multiselect();
            bindSubmit();
            bindOverrideField();
            bindToggle();
            bindShowFraudOverride();
            $('[data-button="toggle"]').parents().next('.hide').toggle()
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

function displayTIDTxnLimits() {
    $.ajax({
        data: {
            override: true,
            siteId: $("#siteID")[0].value,
            profileId: $("#profileID")[0].value,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value,
            TxnLimitData: JSON.stringify(tidLimitData)
        },
        type: 'POST',
        url: '/velocityLimitRow',
        success: function (data) {
            $('#TIDtxnLimitContainer').html(sanitizeHTML(data));
            disableUpdateButtons()
        },
        error: function (data) {
            displayWarningMessage(data.responseText)
        }
    })
}

