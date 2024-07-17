function isClickedInsideConfigModalContent(element) {
    var modalContent = document.getElementById('config-modal-content');
    return modalContent.contains(element);
}

function ishandleClickOutsideConfigModalContent(event) {
    var clickedElement = event.target;
    var modalContent = document.getElementById('config-modal-content');
    if (!isClickedInsideConfigModalContent(clickedElement) && !modalContent.contains(clickedElement) && configModalCheck === 1 && clickedElement.id != "cashback-definitions-edit-modal" && clickedElement.type != "button") {
        emptyToolTipAndCloseConfigModal();
    }
}

document.addEventListener('click', ishandleClickOutsideConfigModalContent);

function renderCashbackTable() {
    modalVars = class {
        static cashbackTableID = '#cashback-table'
        static cashbackSaveID = 'cashback-save'
        static tooltipsClass = '.tooltip'
        static cashbackInputColumnClass = 'cashback-input-column'
        static checkTypeTitle = 'CheckType'
        static checkTypeValues = [1, 2]
        static contentEditableProp = 'contenteditable'
        static domModifiedEvent = 'DOMSubtreeModified'
        static emptyRowText = '12Remove'
        static validate = false
    }

    $('#cashback-loader').remove()

    let tab = $('#cashback')
    let table = $(modalVars.cashbackTableID)
    let jsonTable = new JSONTable(table)
    jsonTable.fromJSON(importJsonFromField());
    let tableColumns = $(modalVars.cashbackTableID + ' th')
    tableColumns.addClass(modalVars.cashbackInputColumnClass)

    let tableCells = $(modalVars.cashbackTableID + ' td')
    tableCells.prop(modalVars.contentEditableProp, true)
    tableCells.on(modalVars.domModifiedEvent, validateCell)
    addRemoveColumn()

    let addButton = $('<button>')
    addButton.prop('id', 'cashback-add-row')
    addButton.prop('class', 'btn btn-secondary')
    addButton.prop('type', 'button')
    addButton.text('Add Row')
    addButton.click(addCashbackConfigRow)
    tab.append(addButton)

    let saveButton = $('<button>')
    saveButton.prop('id', modalVars.cashbackSaveID)
    saveButton.prop('class', 'btn btn-primary float-end')
    saveButton.prop('type', 'button')
    saveButton.text('Update')
    saveButton.click(updateDefinitions)
    tab.append(saveButton)


    setSelectColumn(modalVars.checkTypeTitle, modalVars.checkTypeValues)
    initFieldIds();
    modalVars.validate = true
}

function updateDefinitions() {
    removeLastRowIfEmpty();
    if (!fieldsAreSet()) {
        return;
    }

    removeLastColumn()
    revertSelectColumn(modalVars.checkTypeTitle)

    removeLastRowIfEmpty();
    writeToField(cashbackDefinitionsToJson());

    addRemoveColumn()
    setSelectColumn(modalVars.checkTypeTitle, modalVars.checkTypeValues)

    closeJsonEditModal();
}

/**
 * Checks if all the fields are set, if a field is not set then false is returned and an alert is displayed to the user
 */
function fieldsAreSet() {
    hideAlertDialog();
    let jsonTable = new JSONTable($(modalVars.cashbackTableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return true;
    }

    let aFieldIsEmpty = false;
    json.forEach(function (item, index) {
        if (aFieldIsEmpty) {
            return false;
        }

        let templateErr = "Entry number " + (index + 1) + " has invalid field: ";
        if (!item["BIN"]) {
            templateErr += "BIN"
            aFieldIsEmpty = true;
        } else if (!item["MAX Cashback Amount"]) {
            templateErr += "MAX Cashback Amount"
            aFieldIsEmpty = true;
        } else if (!item["MIN Purchase Amount"]) {
            templateErr += "MIN Purchase Amount"
            aFieldIsEmpty = true;
        }
        showAlertDialog(templateErr + " is Blank")
    });
    return !aFieldIsEmpty;
}

function initFieldIds() {
    let editRows = $(modalVars.cashbackTableID + ' tbody tr');
    editRows.each((i) => {
        let rowTds = editRows[i].getElementsByTagName("td");
        rowTds[0].setAttribute("id", "bin-row-" + i);
        rowTds[1].setAttribute("id", "min-purchase-amount-row-" + i);
        rowTds[2].setAttribute("id", "max-cashback-amount-row-" + i);
        rowTds[3].getElementsByTagName("select")[0].setAttribute("id", "check-type-row-" + i);
        rowTds[3].getElementsByTagName("select")[0].setAttribute("class", "form-select");
        rowTds[4].getElementsByTagName("button")[0].setAttribute("id", "remove-button-row-" + i);
    });
}

function showAlertDialog(message) {
    let alertWindow = $("#cashback-modal-alert-window");
    alertWindow.text(message);
    alertWindow.show();
}

function hideAlertDialog() {
    let alertWindow = $("#cashback-modal-alert-window");
    alertWindow.text("");
    alertWindow.hide();
}

function removeLastRowIfEmpty() {
    let jsonTable = new JSONTable($(modalVars.cashbackTableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return;
    }

    for (let i = 0; i < json.length; i++) {
        let lastRow = json[i];
        if (!lastRow["BIN"] && !lastRow["MAX Cashback Amount"] && !lastRow["MIN Purchase Amount"]) {
            $(modalVars.cashbackTableID + ' tbody tr:last-child').remove();
        }
    }

}

function cashbackDefinitionsToJson() {
    let jsonTable = new JSONTable($(modalVars.cashbackTableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return ""
    } else {
        return JSON.stringify(json);
    }
}

function writeToField(json) {
    $("#" + jsonEditModalOutputFieldId).val(json);
    $('#' + jsonEditModalOutputFieldId).attr('value', json);

}

function importJsonFromField() {
    let fieldVal = $("#" + jsonEditModalOutputFieldId).val();
    if (fieldVal == null || fieldVal === "") {
        return JSON.parse(`[{"BIN":"","MIN Purchase Amount":"","MAX Cashback Amount":"","CheckType":"1"}]`);
    } else {
        return JSON.parse(JSON.stringify(JSON.parse(fieldVal) , ["BIN","MIN Purchase Amount","MAX Cashback Amount","CheckType"]));
    }
}

function addCashbackConfigRow() {
    modalVars.validate = false
    let newRow = $(modalVars.cashbackTableID + ' tr:last').clone()
    let newChildren = newRow.children()
    newChildren.on(modalVars.domModifiedEvent, validateCell)
    newChildren.text('')
    newChildren.removeAttr('aria-describedby')
    newChildren.last().append(newRemoveButton())
    $(modalVars.cashbackTableID + ' tbody').append(newRow)
    revertSelectColumn(modalVars.checkTypeTitle)
    setSelectColumn(modalVars.checkTypeTitle, modalVars.checkTypeValues)
    initFieldIds();
    modalVars.validate = true
}

function removeRow() {
    if ($(this).parent().parent().parent().children().length === 1) {
        $(this).parent().parent().children().text('')
        removeLastColumn()
        addRemoveColumn()
        setSelectColumn(modalVars.checkTypeTitle, modalVars.checkTypeValues)
        return
    }

    $(this).parent().parent().remove()

    tooltipCleanup()
    shouldDisableSaveButton(isAnyTooltips())
    initFieldIds();
}

function removeLastColumn() {
    $(modalVars.cashbackTableID + ' th:last-child, ' + modalVars.cashbackTableID + ' td:last-child').remove();
}

function addRemoveColumn() {
    let newColumn = $('<td>')
    newColumn.append(newRemoveButton())

    $(modalVars.cashbackTableID + ' tbody tr').append(newColumn)
    $(modalVars.cashbackTableID + ' thead tr').append('<th class="remove-column">')
}

function newRemoveButton() {
    let removeButton = $('<button>')
    removeButton.prop('class', 'btn btn-secondary')
    removeButton.prop('type', 'button')
    removeButton.click(removeRow)
    return removeButton.text('Remove')
}

function validateCell() {
    if (!modalVars.validate) {
        return
    }
    let cell = $(this)
    let title = cell.parents('table').find('th').eq(cell.index()).text();
    let config = configForColumn(title)
    // limit cell size to maxLength characters
    if (cell.text().length > config.maxLength) {
        cell.text(cell.text().substring(0, config.maxLength))

        setCaratPosition(this, config.maxLength)
    }

    let valid = RegExp(config.regex).test(cell.text())

    if (emptyCashbackConfig()) {
        // remove all cells with a tooltip
        $(modalVars.cashbackTableID + ' td[aria-describedby*="tooltip"]').tooltip('dispose')
    } else if (valid) {
        cell.tooltip('dispose')
    } else {
        cell.tooltip({ title: config.error, trigger: 'manual', animation: false }).tooltip('show')
    }

    tooltipCleanup()
    // disable save button if there are any tooltips
    shouldDisableSaveButton(isAnyTooltips() || emptyCashbackConfig())
}

// a side effect of setting cell.text the carat is moved to the beginning. This moves it arbitrarily
function setCaratPosition(cell, location) {
    let range = document.createRange()
    let selection = window.getSelection()

    range.setStart(cell.childNodes[0], location)
    range.collapse(true)

    selection.removeAllRanges()
    selection.addRange(range)

    cell.focus()
}

// returns true if there is only a single line and all fields are empty
function emptyCashbackConfig() {
    return $(modalVars.cashbackTableID + ' tbody tr').length === 1 && $(modalVars.cashbackTableID + ' tr td').text() === modalVars.emptyRowText
}

function configForColumn(title) {
    switch (title) {
        case 'BIN':
            return {
                regex: '^[0-9]+(>[0-9]+)?$',
                error: 'Invalid BIN range. Format must be [number] or [number]>[number]',
                maxLength: 40
            }
        case 'MIN Purchase Amount':
        case 'MAX Cashback Amount':
            return {
                regex: '^[0-9]+$',
                error: 'Invalid amount, can only be a number',
                maxLength: 15
            }
        default:
            return {
                regex: '',
                error: '',
                maxLength: 15
            }
    }
}

// bootstrap tooltips sometimes create stray tooltips, this removes all invalid tooltips
function tooltipCleanup() {
    for (let tooltip of $(modalVars.tooltipsClass)) {
        if ($('[aria-describedby=' + tooltip.id + ']').length === 0) {
            tooltip.remove()
        }
    }
}

function shouldDisableSaveButton(shouldDisable) {
    $('#' + modalVars.cashbackSaveID).prop('disabled', shouldDisable)
}

function isAnyTooltips() {
    return $(modalVars.tooltipsClass).length !== 0
}

function setSelectColumn(title, values) {
    let columnHeader = $('th:contains(' + title + ')')
    columnHeader.removeClass(modalVars.cashbackInputColumnClass)
    let index = columnHeader.index() + 1

    $(modalVars.cashbackTableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.prop(modalVars.contentEditableProp, false)

        let selector = $('<select>')

        for (let val of values) {
            selector.append($('<option>').prop('value', val).text(val))
        }

        // currently a switch so if more values are added in the future, they can without converting if to switch
        switch (cell.text()) {
            case '2':
                selector.find('option[value="2"]').prop('selected', true)
        }

        cell.html(selector)
    })
}

function revertSelectColumn(title) {
    let index = $('th:contains(' + title + ')').index() + 1

    $(modalVars.cashbackTableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.text(cell.find('option:selected').val())
    })
}

function emptyToolTipAndCloseConfigModal() {
    $(modalVars.cashbackTableID + ' td[aria-describedby*="tooltip"]').hide();
    $("#config-edit-modal").modal('hide');
    configModalCheck = 0

}