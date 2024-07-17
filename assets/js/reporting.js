let type = "week";
// A simple counter to track the number of data sets currently retrieved
let chartsRetrieved = -1;
// Tracks whether the filter request has been cancelled or not
let activeFilterRequest;

$(document).ready(function() {
    const dateTimeStep = 5

    $("#after").datetimepicker({
        id: 'after-datetime',
        step: dateTimeStep
    })

    $("#before").datetimepicker({
        id: 'before-datetime',
        step: dateTimeStep
    })

    const acquirersMultiSelect = $('#acquirersMulti')
    acquirersMultiSelect.multiselect({
        maxHeight: 200,
        numberDisplayed: 1,
        onInitialized: (select) => {
            select.next().children('button').attr('id', 'acquirersMulti-dropdown')
        }
    });

    $("#clear-filters").click(() => {
        $('#before')[0].value = "";
        $('#after')[0].value = "";
        acquirersMultiSelect.multiselect('deselectAll', false);
        acquirersMultiSelect.multiselect('updateButtonText');
    });
})

function filterCharts() {
    resetCharts();
    posMIDs();
}

function resetCharts() {
    $("#total-mids").html("0");
    $("#total-tids").html("0");
    $("#transacting-mids").html("0");
    $("#transacting-tids").html("0");
    $("#total-transactions").html("0");
    $("#total-transaction-value").html("0");
    $("#approved-transactions").html("0");
    $("#declined-transactions").html("0");
    $("#Total-Volume").disposeFusionCharts();
    $("#Approved-Transaction-Volume").disposeFusionCharts();
    $("#Declined-Transaction-Volume").disposeFusionCharts();
    $("#Top-10-Decline-Reasons").disposeFusionCharts();
    $("#Top-10-MIDs").disposeFusionCharts();
    $("#Top-10-TIDs").disposeFusionCharts();
    $("#Top-10-Card-Types").disposeFusionCharts();
    hideWarning();
}

function filterChartsCustom() {
    resetCharts();
    startLoadingCharts();
    filterCharts();
}

function posMIDs() {
    var data = {
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    displayCurrentFetch("Total MIDs");

    activeFilterRequest = $.ajax({
        url: "/reporting/totalMIDs",
        method: "POST",
        data: data,
        success: function (d) {
            $("#total-mids").html(d);
            posTIDs();
        },
        error: splitErrorFromAbort
    });
}

function posTIDs() {
    var data = {
        csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
    }
    displayCurrentFetch("Total TIDs");

    activeFilterRequest = $.ajax({
        url: "/reporting/totalTIDs",
        data: data,
        method: "POST",
        data: data,
        success: function (d) {
            $("#total-tids").html(d);
            TransactingMIDs();
        },
        error: splitErrorFromAbort
    });
}

function TransactingMIDs() {
    displayCurrentFetch("Transacting MIDs");
    activeFilterRequest = $.ajax({
        url: "/reporting/TransactingMIDs",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#transacting-mids").html(d);
            TransactingTIDs();
        },
        error: splitErrorFromAbort
    });
}

function TransactingTIDs() {
    displayCurrentFetch("Transacting TIDs");
    activeFilterRequest = $.ajax({
        url: "/reporting/TransactingTIDs",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#transacting-tids").html(d);
            TotalTransactions();
        },
        error: splitErrorFromAbort
    });
}

function TotalTransactions() {
    displayCurrentFetch("Total Transactions");
    activeFilterRequest = $.ajax({
        url: "/reporting/TotalTransactions",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#total-transactions").html(d);
            TotalTransactionValue();
        },
        error: splitErrorFromAbort
    });
}

function TotalTransactionValue() {
    displayCurrentFetch("Total Transaction Value");
    activeFilterRequest = $.ajax({
        url: "/reporting/TotalTransactionValue",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#total-transaction-value").html(d);
            ApprovedTransactions();
        },
        error: splitErrorFromAbort
    });
}

function ApprovedTransactions() {
    displayCurrentFetch("Approved Transactions");
    activeFilterRequest = $.ajax({
        url: "/reporting/ApprovedTransactions",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#approved-transactions").html(d);
            DeclinedTransactions();
        },
        error: splitErrorFromAbort
    });
}

function DeclinedTransactions() {
    displayCurrentFetch("Declined Transactions");
    activeFilterRequest = $.ajax({
        url: "/reporting/DeclinedTransactions",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#declined-transactions").html(d);
            drawTotalVolume();
        },
        error: splitErrorFromAbort
    });
}

function drawTotalVolume() {
    const chartTitle = "Transaction Count";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTransactionVolume",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Total-Volume").insertFusionCharts({
                id: 'total-volume-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "Date",
                        "yAxisName": "Count",
                        "canvasPadding": "30",
                        "theme": "candy",
                    },
                    "data": JSON.parse(d)
                }
            });
            drawApprovedTransactionVolume();
        },
        error: splitErrorFromAbort
    });
}

function drawApprovedTransactionVolume() {
    const chartTitle = "Approved Transaction Count";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawApprovedTransactionVolume",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Approved-Transaction-Volume").insertFusionCharts({
                id: 'approved-txns-volume-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Approved Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "Date",
                        "yAxisName": "Count",
                        "canvasPadding": "30",
                        "theme": "candy",
                    },
                    "data": JSON.parse(d)
                }
            });
            drawDeclinedTransactionVolume();
        },
        error: splitErrorFromAbort
    });
}

function drawDeclinedTransactionVolume() {
    const chartTitle = "Declined Transaction Count";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawDeclinedTransactionVolume",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Declined-Transaction-Volume").insertFusionCharts({
                id: 'declined-txns-volume-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Declined Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "Date",
                        "yAxisName": "Count",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTotalValue();
        },
        error: splitErrorFromAbort
    });
}

function drawTotalValue() {
    const chartTitle = "Transaction Value";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTransactionValue",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Total-Value").insertFusionCharts({
                id: 'total-value-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "Date",
                        "yAxisName": "Amount (AED)",
                        "canvasPadding": "30",
                        "theme": "candy",
                    },
                    "data": JSON.parse(d)
                }
            });
            drawApprovedTransactionValue();
        },
        error: splitErrorFromAbort
    });
}

function drawApprovedTransactionValue() {
    const chartTitle = "Approved Transaction Value";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawApprovedTransactionValue",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Approved-Transaction-Value").insertFusionCharts({
                id: 'approved-txn-value-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Approved Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "Date",
                        "yAxisName": "Amount (AED)",
                        "canvasPadding": "30",
                        "theme": "candy",
                    },
                    "data": JSON.parse(d)
                }
            });
            drawDeclinedTransactionValue();
        },
        error: splitErrorFromAbort
    });
}

function drawDeclinedTransactionValue() {
    const chartTitle = "Declined Transaction Value";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawDeclinedTransactionValue",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Declined-Transaction-Value").insertFusionCharts({
                id: 'declined-txn-value-chart',
                type: "Line",
                width: "30%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Declined Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "Date",
                        "yAxisName": "Amount (AED)",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10DeclineReasons();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10DeclineReasons() {
    const chartTitle = "Top 10 Decline Reasons";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10DeclineReasons",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-Decline-Reasons").insertFusionCharts({
                id: 'top-ten-decline-reasons-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Declined Transactions in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "Reason",
                        "yAxisName": "Total",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10MIDs();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10MIDs() {
    const chartTitle = "Top 10 Transacting MIDs";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10TransactingMIDs",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-MIDs").insertFusionCharts({
                id: 'top-ten-mids-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transacting MIDs in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "MID",
                        "yAxisName": "Total",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10TIDs();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10TIDs() {
    const chartTitle = "Top 10 Transacting TIDs";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10TransactingTIDs",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-TIDs").insertFusionCharts({
                id: 'top-ten-tids-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transacting TIDs in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "TID",
                        "yAxisName": "Total",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10CardTypes();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10CardTypes() {
    const chartTitle = "Top 10 Card Types";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10CardTypes",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-Card-Types").insertFusionCharts({
                id: 'top-ten-card-types-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Card payments taken in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "xAxisName": "Card Type",
                        "yAxisName": "Total",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10MIDValues();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10MIDValues() {
    const chartTitle = "Top 10 Transacting MIDs";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10TransactingMIDValues",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-MID-Values").insertFusionCharts({
                id: 'top-ten-mid-values-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transacting MIDs in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "MID",
                        "yAxisName": "Amount (AED)",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10TIDValues();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10TIDValues() {
    const chartTitle = "Top 10 Transacting TIDs";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10TransactingTIDValues",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-TID-Values").insertFusionCharts({
                id: 'top-ten-tid-values-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Transacting TIDs in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "TID",
                        "yAxisName": "Amount (AED)",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            drawTop10CardTypeValues();
        },
        error: splitErrorFromAbort
    });
}

function drawTop10CardTypeValues() {
    const chartTitle = "Top 10 Card Types";
    displayCurrentFetch(chartTitle);
    activeFilterRequest = $.ajax({
        url: "/reporting/drawTop10CardTypeValues",
        method: "POST",
        data: {
            Before: $("#before")[0].value,
            After: $("#after")[0].value,
            Acquirers: $("#acquirersMulti").val(),
            Type: type,
            csrfmiddlewaretoken: $("input[name=csrfmiddlewaretoken]")[0].value
        },
        success: function (d) {
            $("#Top-10-Card-Type-Values").insertFusionCharts({
                id: 'top-ten-card-type-values-chart',
                type: "bar3d",
                width: "20%",
                height: "400",
                dataFormat: "json",
                dataEmptyMessage: "No Card payments taken in this period",
                dataSource: {
                    // Chart Configuration
                    "chart": {
                        "caption": chartTitle,
                        "subCaption": "(AED)",
                        "xAxisName": "Card Type",
                        "yAxisName": "Amount (AED)",
                        "theme": "candy",
                    },
                    // Chart Data
                    "data": JSON.parse(d)
                }
            });
            finishLoadingCharts();
        },
        error: splitErrorFromAbort
    });
}

// Provides user with quantifiable progress
function displayCurrentFetch(chartName) {
    chartsRetrieved ++;
    $("#data-being-fetched").text("Fetching..." + chartName);
    $("#data-successfully-retrieved").text(chartsRetrieved + "/21 retrieved");
}

// Shows the loading modal
function startLoadingCharts() {
    chartsRetrieved = -1;
    $("#reporting-charts").hide();
    $("#report-loading-modal").show();
}

// Hides the loading modal and shows charts
function finishLoadingCharts() {
    $("#report-loading-modal").hide();
    $("#reporting-charts").show();
}

// Cancel loading charts
function cancelGenerateReport() {
    $("#reporting-charts").hide();
    $("#report-loading-modal").hide();
    activeFilterRequest.abort();
    resetCharts();
}

// Calling .abort() on an Ajax request does not stop the server side of the request. When the server side request
// completes it will throw an error on the calling Ajax. This function determines if the error thrown is genuine or
// simply the aborted response
function splitErrorFromAbort(data) {
    if (data.statusText !== "abort") {
        finishLoadingCharts();
        displayWarningMessage(data.responseText);
    }
}

function setActive(e) {
    type = e;

    $("#monthFilter").removeClass("active");
    $("#weekFilter").removeClass("active");
    $("#dayFilter").removeClass("active");
    $("#customFilter").removeClass("active");

    if (type == "month") {
        $("#monthFilter").addClass("active");
    } else if (type == "week" ) {
        $("#weekFilter").addClass("active");
    } else if (type == "day") {
        $("#dayFilter").addClass("active");
    } else if (type == "custom") {
        $("#customFilter").addClass("active");
    }
}