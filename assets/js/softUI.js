function isClickedInsideSoftUIModalContent(element) {
    var modalContent = document.getElementById('soft-ui-modal-content');
    return modalContent.contains(element);
}

function ishandleClickOutsideSoftUIModalContent(event) {
    var clickedElement = event.target;
    var modalContent = document.getElementById('soft-ui-modal-content');
    if (!isClickedInsideSoftUIModalContent(clickedElement) && !modalContent.contains(clickedElement) && softUIModalCheck === 1 && clickedElement.id != "soft-ui-mccDetails-edit-modal" && clickedElement.type != "button") {
        emptyToolTipAndCloseSoftUIModal();
    }
}

document.addEventListener('click', ishandleClickOutsideSoftUIModalContent);

function renderSoftUITable(options) {
    var separatedArray = [];
    for (var i = 0; i < options.length; i++) {
        separatedArray.push(options[i].Option)
    };

    modalSoftUIVars = class {
        static softUITableID = '#soft-ui-table'
        static softUISaveID = 'soft-ui-save'
        static tooltipsSoftUIClass = '.tooltip'
        static softUIInputColumnClass = 'soft-ui-input-column'
        static checkTypeTitle = 'APM'
        static checkTypeValues = separatedArray
        static contentSoftUIEditableProp = 'contenteditable'
        static domSoftUIModifiedEvent = 'DOMSubtreeModified'
        static emptySoftUIRowText = '12Remove'
        static validateSoftUI = false
    }

    $('#soft-ui-loader').remove()

    let tab = $('#softUI')
    let softUITable = $(modalSoftUIVars.softUITableID)
    let jsonTable = new JSONTable(softUITable)
    jsonTable.fromJSON(importJsonFromSoftUIField());
    let tableColumns = $(modalSoftUIVars.softUITableID + ' th')
    tableColumns.addClass(modalSoftUIVars.softUIInputColumnClass)

    let tableCells = $(modalSoftUIVars.softUITableID + ' td')
    tableCells.prop(modalSoftUIVars.contentSoftUIEditableProp, true)
    tableCells.on(modalSoftUIVars.domSoftUIModifiedEvent, validateSoftUICell)
    addRemoveSoftUIColumn()

    let addButton = $('<button>')
    addButton.prop('id', 'soft-ui-add-row')
    addButton.prop('class', 'btn btn-secondary')
    addButton.prop('type', 'button')
    addButton.text('Add Row')
    addButton.click(addSoftUIConfigRow)
    tab.append(addButton)

    let saveButton = $('<button>')
    saveButton.prop('id', modalSoftUIVars.softUISaveID)
    saveButton.prop('class', 'btn btn-primary float-end')
    saveButton.prop('type', 'button')
    saveButton.text('Update')
    saveButton.click(updateSoftUIDefinitions)
    tab.append(saveButton)

    setSelectSoftUIColumn(modalSoftUIVars.checkTypeTitle, modalSoftUIVars.checkTypeValues)
    initSoftUIFieldIds();
    modalSoftUIVars.validateSoftUI = true
}

function OnChangeSelect(value){
    $('#' + modalSoftUIVars.softUISaveID).prop('disabled', false)
}

function updateSoftUIDefinitions() {
    let duplicate = checkAPMValidation();

    if (duplicate) {
        showAlertSoftUIDialog("Multiple selection of APM is not allowed")
        return;
    }

    removeLastSoftUIRowIfEmpty();
    if (!fieldsSoftUIAreSet()) {
        return;
    }

    removeLastSoftUIColumn()
    revertSelectSoftUIColumn(modalSoftUIVars.checkTypeTitle)

    removeLastSoftUIRowIfEmpty();
    writeToSoftUIField(softUIDefinitionsToJson());

    addRemoveSoftUIColumn()
    setSelectSoftUIColumn(modalSoftUIVars.checkTypeTitle, modalSoftUIVars.checkTypeValues)

    closeJsonEditSoftUIModal();
}

/**
 * Checks if all the fields are set, if a field is not set then false is returned and an alert is displayed to the user
 */
function fieldsSoftUIAreSet() {
    hideAlertSoftUIDialog();
    let jsonTable = new JSONTable($(modalSoftUIVars.softUITableID));
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
        if (!item["APM"]) {
            templateErr += "APM"
            aFieldIsEmpty = true;
        } else if (!item["MCC"]) {
            templateErr += "MCC"
            aFieldIsEmpty = true;
        }
        showAlertSoftUIDialog(templateErr + " is Blank")
    });
    return !aFieldIsEmpty;
}

function initSoftUIFieldIds() {
    let editRows = $(modalSoftUIVars.softUITableID + ' tbody tr');
    editRows.each((i) => {
        let rowTds = editRows[i].getElementsByTagName("td");
        rowTds[0].getElementsByTagName("select")[0].setAttribute("id", "apm-row-" + i);
        rowTds[0].getElementsByTagName("select")[0].setAttribute("class", "form-select");
        rowTds[0].getElementsByTagName("select")[0].setAttribute("style", "width: 50%;");
        rowTds[0].getElementsByTagName("select")[0].setAttribute("onchange", "OnChangeSelect(this.value)");
        rowTds[1].setAttribute("id", "mcc-row-" + i);
        rowTds[2].getElementsByTagName("button")[0].setAttribute("id", "remove-button-row-" + i);
    });
}

function showAlertSoftUIDialog(message) {
    let alertWindow = $("#soft-ui-modal-alert-window");
    alertWindow.text(message);
    alertWindow.show();
}

function hideAlertSoftUIDialog() {
    let alertWindow = $("#soft-ui-modal-alert-window");
    alertWindow.text("");
    alertWindow.hide();
}

function removeLastSoftUIRowIfEmpty() {
    let jsonTable = new JSONTable($(modalSoftUIVars.softUITableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return;
    }

    for (let i = 0; i < json.length; i++) {
        let lastRow = json[i];
        if (!lastRow["APM"] && !lastRow["MCC"]) {
            $(modalSoftUIVars.softUITableID + ' tbody tr:last-child').remove();
        }
    }
}

function softUIDefinitionsToJson() {
    let jsonTable = new JSONTable($(modalSoftUIVars.softUITableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return ""
    } else {
        return JSON.stringify(json);
    }
}

function writeToSoftUIField(json) {
    $("#" + jsonEditModalSoftUIOutputFieldId).val(json);
    $('#' + jsonEditModalSoftUIOutputFieldId).attr('value', json);
}

function importJsonFromSoftUIField() {
    let fieldVal = $("#" + jsonEditModalSoftUIOutputFieldId).val();
    if (fieldVal == null || fieldVal === "") {
        return JSON.parse(`[{"APM":"Default", "MCC": ""}]`);
    } else {
        return JSON.parse(fieldVal);
    }
}

function addSoftUIConfigRow() {
    modalSoftUIVars.validateSoftUI = false
    let newRow = $(modalSoftUIVars.softUITableID + ' tr:last').clone()
    let newChildren = newRow.children()
    newChildren.on(modalSoftUIVars.domSoftUIModifiedEvent, validateSoftUICell)
    newChildren.text('')
    newChildren.removeAttr('aria-describedby')
    newChildren.last().append(newRemoveSoftUIButton())
    $(modalSoftUIVars.softUITableID + ' tbody').append(newRow)
    revertSelectSoftUIColumn(modalSoftUIVars.checkTypeTitle)
    setSelectSoftUIColumn(modalSoftUIVars.checkTypeTitle, modalSoftUIVars.checkTypeValues)
    initSoftUIFieldIds();
    modalSoftUIVars.validateSoftUI = true
}

function removeSoftUIRow() {
    if ($(this).parent().parent().parent().children().length === 1) {
        $(this).parent().parent().children('td#mcc-row-0').html('')
        return
    }

    $(this).parent().parent().remove()

    tooltipSoftUICleanup()
    shouldDisableSoftUISaveButton(isAnySoftUITooltips())
    initSoftUIFieldIds();
}

function removeLastSoftUIColumn() {
    $(modalSoftUIVars.softUITableID + ' th:last-child, ' + modalSoftUIVars.softUITableID + ' td:last-child').remove();
}

function addRemoveSoftUIColumn() {
    let newColumn = $('<td>')
    newColumn.append(newRemoveSoftUIButton())

    $(modalSoftUIVars.softUITableID + ' tbody tr').append(newColumn)
    $(modalSoftUIVars.softUITableID + ' thead tr').append('<th class="remove-column">')
}

function newRemoveSoftUIButton() {
    let removeButton = $('<button>')
    removeButton.prop('class', 'btn btn-secondary')
    removeButton.prop('type', 'button')
    removeButton.click(removeSoftUIRow)
    return removeButton.text('Remove')
}

function validateSoftUICell() {
    if (!modalSoftUIVars.validateSoftUI) {
        return
    }
    let cell = $(this)
    let title = cell.parents('table').find('th').eq(cell.index()).text();
    let config = configForSoftUIColumn(title)

    // limit cell size to maxLength characters
    if (cell.text().length > config.maxLength) {
        cell.text(cell.text().substring(0, config.maxLength))

        setSoftUICaratPosition(this, config.maxLength)
    }

    let valid = RegExp(config.regex).test(cell.text())

    if (emptySoftUIConfig()) {
        // remove all cells with a tooltip
        $(modalSoftUIVars.softUITableID + ' td[aria-describedby*="tooltip"]').tooltip('dispose')
    } else if (valid) {
        cell.tooltip('dispose')
    } else {
        cell.tooltip({ title: config.error, trigger: 'manual', animation: false }).tooltip('show')
    }

    tooltipSoftUICleanup()

    // disable save button if there are any tooltips
    shouldDisableSoftUISaveButton(isAnySoftUITooltips() || emptySoftUIConfig())
}

function shouldDisableSoftUISaveButton(shouldDisable) {
    $('#' + modalSoftUIVars.softUISaveID).prop('disabled', shouldDisable)
}

function isAnySoftUITooltips() {
    return $(modalSoftUIVars.tooltipsSoftUIClass).length !== 0
}

// a side effect of setting cell.text the carat is moved to the beginning. This moves it arbitrarily
function setSoftUICaratPosition(cell, location) {
    let range = document.createRange()
    let selection = window.getSelection()

    range.setStart(cell.childNodes[0], location)
    range.collapse(true)

    selection.removeAllRanges()
    selection.addRange(range)

    cell.focus()
}

// returns true if there is only a single line and all fields are empty
function emptySoftUIConfig() {
    return $(modalSoftUIVars.softUITableID + ' tbody tr').length === 1 && $(modalSoftUIVars.softUITableID + ' tr td').text() === modalSoftUIVars.emptySoftUIRowText
}

function configForSoftUIColumn(title) {
    switch (title) {
        case 'MCC':
            return {
                regex: '^[0-9]+$',
                error: 'Invalid MCC, can only be a number',
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
function tooltipSoftUICleanup() {
    for (let tooltip of $(modalSoftUIVars.tooltipsSoftUIClass)) {
        if ($('[aria-describedby=' + tooltip.id + ']').length === 0) {
            tooltip.remove()
        }
    }
}

function setSelectSoftUIColumn(title, values) {
    let columnHeader = $('th:contains(' + title + ')')
    columnHeader.removeClass(modalSoftUIVars.softUIInputColumnClass)
    let index = columnHeader.index() + 1

    $(modalSoftUIVars.softUITableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.prop(modalSoftUIVars.contentSoftUIEditableProp, false)

        let selector = $('<select>')

        for (let val of values) {
            const newOption = $('<option>', {
                value: val,
                text: val,
            });
            selector.append(newOption)
        }

        selector.find('option[value="'+cell.text()+'"]').prop('selected', true)

        cell.html(selector)
    })
}

function revertSelectSoftUIColumn(title) {
    let index = $('th:contains(' + title + ')').index() + 1

    $(modalSoftUIVars.softUITableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.text(cell.find('option:selected').val())
    })
}

function emptyToolTipAndCloseSoftUIModal() {
    $(modalSoftUIVars.softUITableID + ' td[aria-describedby*="tooltip"]').hide();
    $("#soft-ui-config-edit-modal").modal('hide');
    softUIModalCheck = 0

}

function checkAPMValidation() {
    let selectElements = document.querySelectorAll('table'+modalSoftUIVars.softUITableID+' select');
    let selectedValues = [];
    let duplicate = false;
    selectElements.forEach(select => {
        let selectedValue = select.value;
        if (selectedValues.includes(selectedValue)) {
            duplicate = true;
        } else {
            selectedValues.push(selectedValue);
        }
    });
    return duplicate;
}