function isClickedInsideDpoModalContent(element) {
    var modalContent = document.getElementById('dpo-modal-content');
    return modalContent.contains(element);
}

function ishandleClickOutsideDpoModalContent(event) {
    var clickedElement = event.target;
    var modalContent = document.getElementById('dpo-modal-content');
    if (!isClickedInsideDpoModalContent(clickedElement) && !modalContent.contains(clickedElement) && dpoModalCheck === 1 && clickedElement.id != "dpoMomo-countryDetails-edit-modal" && clickedElement.type != "button") {
        emptyToolTipAndCloseDpoModal();
    }
}

document.addEventListener('click', ishandleClickOutsideDpoModalContent);

function renderDpoTable() {
    modalDpoVars = class {
        static dpoTableID = '#dpo-table'
        static dpoSaveID = 'dpo-save'
        static tooltipsDpoClass = '.tooltip'
        static dpoInputColumnClass = 'dpo-input-column'
        static checkTypeTitle = 'IsDefault'
        static checkTypeValues = [true, false]
        static contentDpoEditableProp = 'contenteditable'
        static domDpoModifiedEvent = 'DOMSubtreeModified'
        static emptyDpoRowText = '12Remove'
        static validateDpo = false
    }

    $('#dpo-loader').remove()

    let tab = $('#dpo')
    let dpoTable = $(modalDpoVars.dpoTableID)
    let jsonTable = new JSONTable(dpoTable)
    jsonTable.fromJSON(importJsonFromDpoField());
    let tableColumns = $(modalDpoVars.dpoTableID + ' th')
    tableColumns.addClass(modalDpoVars.dpoInputColumnClass)

    let tableCells = $(modalDpoVars.dpoTableID + ' td')
    tableCells.prop(modalDpoVars.contentDpoEditableProp, true)
    tableCells.on(modalDpoVars.domDpoModifiedEvent, validateDpoCell)
    addRemoveDpoColumn()

    let addButton = $('<button>')
    addButton.prop('id', 'dpo-add-row')
    addButton.prop('class', 'btn btn-secondary')
    addButton.prop('type', 'button')
    addButton.text('Add Row')
    addButton.click(addDpoConfigRow)
    tab.append(addButton)

    let saveButton = $('<button>')
    saveButton.prop('id', modalDpoVars.dpoSaveID)
    saveButton.prop('class', 'btn btn-primary float-end')
    saveButton.prop('type', 'button')
    saveButton.text('Update')
    saveButton.click(updateDpoDefinitions)
    tab.append(saveButton)

    setSelectDpoColumn(modalDpoVars.checkTypeTitle, modalDpoVars.checkTypeValues)
    initDpoFieldIds();
    modalDpoVars.validateDpo = true
}

function updateDpoDefinitions() {
    let [checkDefault, count] = checkDefaultValidation();

    if (count > 1) {
        showAlertDpoDialog("Please select only one default CountryCode")
        return;
    }

    if (!checkDefault) {
        showAlertDpoDialog("Please select a default CountryCode")
        return;
    }
    removeLastDpoRowIfEmpty();
    if (!fieldsDpoAreSet()) {
        return;
    }

    removeLastDpoColumn()
    revertSelectDpoColumn(modalDpoVars.checkTypeTitle)

    removeLastDpoRowIfEmpty();
    writeToDpoField(dpoDefinitionsToJson());

    addRemoveDpoColumn()
    setSelectDpoColumn(modalDpoVars.checkTypeTitle, modalDpoVars.checkTypeValues)

    closeJsonEditDpoModal();
}

/**
 * Checks if all the fields are set, if a field is not set then false is returned and an alert is displayed to the user
 */
function fieldsDpoAreSet() {
    hideAlertDpoDialog();
    let jsonTable = new JSONTable($(modalDpoVars.dpoTableID));
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
        if (!item["CountryName"]) {
            templateErr += "CountryName"
            aFieldIsEmpty = true;
        } else if (!item["CountryCode"]) {
            templateErr += "CountryCode"
            aFieldIsEmpty = true;
        }
        showAlertDpoDialog(templateErr + " is Blank")
    });
    return !aFieldIsEmpty;
}

function initDpoFieldIds() {
    let editRows = $(modalDpoVars.dpoTableID + ' tbody tr');
    editRows.each((i) => {
        let rowTds = editRows[i].getElementsByTagName("td");
        rowTds[0].setAttribute("id", "country-name-row-" + i);
        rowTds[1].setAttribute("id", "country-code-row-" + i);
        rowTds[2].getElementsByTagName("select")[0].setAttribute("id", "country-default-row-" + i);
        rowTds[2].getElementsByTagName("select")[0].setAttribute("class", "form-select");
        rowTds[2].getElementsByTagName("select")[0].setAttribute("style", "width: 29%;");
        rowTds[3].getElementsByTagName("button")[0].setAttribute("id", "remove-button-row-" + i);
    });
}

function showAlertDpoDialog(message) {
    let alertWindow = $("#dpo-modal-alert-window");
    alertWindow.text(message);
    alertWindow.show();
}

function hideAlertDpoDialog() {
    let alertWindow = $("#dpo-modal-alert-window");
    alertWindow.text("");
    alertWindow.hide();
}

function removeLastDpoRowIfEmpty() {
    let jsonTable = new JSONTable($(modalDpoVars.dpoTableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return;
    }

    for (let i = 0; i < json.length; i++) {
        let lastRow = json[i];
        if (!lastRow["CountryName"] && !lastRow["CountryCode"]) {
            $(modalDpoVars.dpoTableID + ' tbody tr:last-child').remove();
        }
    }
}

function dpoDefinitionsToJson() {
    let jsonTable = new JSONTable($(modalDpoVars.dpoTableID));
    let json = jsonTable.toJSON();
    if (json.length === 0) {
        return ""
    } else {
        return JSON.stringify(json);
    }
}

function writeToDpoField(json) {
    $("#" + jsonEditModalDpoOutputFieldId).val(json);
    $('#' + jsonEditModalDpoOutputFieldId).attr('value', json);
}

function importJsonFromDpoField() {
    let fieldVal = $("#" + jsonEditModalDpoOutputFieldId).val();
    if (fieldVal == null || fieldVal === "") {
        return JSON.parse(`[{"CountryName":"","CountryCode":"", "IsDefault": "false"}]`);
    } else {
        return JSON.parse(fieldVal);
    }
}

function addDpoConfigRow() {
    modalDpoVars.validateDpo = false
    let newRow = $(modalDpoVars.dpoTableID + ' tr:last').clone()
    let newChildren = newRow.children()
    newChildren.on(modalDpoVars.domDpoModifiedEvent, validateDpoCell)
    newChildren.text('')
    newChildren.removeAttr('aria-describedby')
    newChildren.last().append(newRemoveDpoButton())
    $(modalDpoVars.dpoTableID + ' tbody').append(newRow)
    revertSelectDpoColumn(modalDpoVars.checkTypeTitle)
    setSelectDpoColumn(modalDpoVars.checkTypeTitle, modalDpoVars.checkTypeValues)
    initDpoFieldIds();
    modalDpoVars.validateDpo = true
}

function removeDpoRow() {
    if ($(this).parent().parent().parent().children().length === 1) {
        $(this).parent().parent().children().text('')
        removeLastDpoColumn()
        addRemoveDpoColumn()
        setSelectDpoColumn(modalDpoVars.checkTypeTitle, modalDpoVars.checkTypeValues)
        return
    }

    $(this).parent().parent().remove()

    tooltipDpoCleanup()
    shouldDisableDpoSaveButton(isAnyDpoTooltips())
    initDpoFieldIds();
}

function removeLastDpoColumn() {
    $(modalDpoVars.dpoTableID + ' th:last-child, ' + modalDpoVars.dpoTableID + ' td:last-child').remove();
}

function addRemoveDpoColumn() {
    let newColumn = $('<td>')
    newColumn.append(newRemoveDpoButton())

    $(modalDpoVars.dpoTableID + ' tbody tr').append(newColumn)
    $(modalDpoVars.dpoTableID + ' thead tr').append('<th class="remove-column">')
}

function newRemoveDpoButton() {
    let removeButton = $('<button>')
    removeButton.prop('class', 'btn btn-secondary')
    removeButton.prop('type', 'button')
    removeButton.click(removeDpoRow)
    return removeButton.text('Remove')
}

function validateDpoCell() {
    if (!modalDpoVars.validateDpo) {
        return
    }
    let cell = $(this)
    let title = cell.parents('table').find('th').eq(cell.index()).text();
    let config = configForDpoColumn(title)

    // limit cell size to maxLength characters
    if (cell.text().length > config.maxLength) {
        cell.text(cell.text().substring(0, config.maxLength))

        setDpoCaratPosition(this, config.maxLength)
    }

    let valid = RegExp(config.regex).test(cell.text())

    if (emptyDpoConfig()) {
        // remove all cells with a tooltip
        $(modalDpoVars.dpoTableID + ' td[aria-describedby*="tooltip"]').tooltip('dispose')
    } else if (valid) {
        cell.tooltip('dispose')
    } else {
        cell.tooltip({ title: config.error, trigger: 'manual', animation: false }).tooltip('show')
    }

    tooltipDpoCleanup()
    // disable save button if there are any tooltips
    shouldDisableDpoSaveButton(isAnyDpoTooltips() || emptyDpoConfig())
}

// a side effect of setting cell.text the carat is moved to the beginning. This moves it arbitrarily
function setDpoCaratPosition(cell, location) {
    let range = document.createRange()
    let selection = window.getSelection()

    range.setStart(cell.childNodes[0], location)
    range.collapse(true)

    selection.removeAllRanges()
    selection.addRange(range)

    cell.focus()
}

// returns true if there is only a single line and all fields are empty
function emptyDpoConfig() {
    return $(modalDpoVars.dpoTableID + ' tbody tr').length === 1 && $(modalDpoVars.dpoTableID + ' tr td').text() === modalDpoVars.emptyDpoRowText
}

function configForDpoColumn(title) {
    switch (title) {
        case 'CountryName':
            return {
                regex: "^(?!\\s)[a-zA-Z\\s\\\\\\/\\-\\)\\(\\`\\.\\\\\"\\']+(?<!\\s)$",
                error: 'Invalid CountryName.',
                maxLength: 40
            }
        case 'CountryCode':
            return {
                regex: '^[0-9]+$',
                error: 'Invalid CountryCode, can only be a number',
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
function tooltipDpoCleanup() {
    for (let tooltip of $(modalDpoVars.tooltipsDpoClass)) {
        if ($('[aria-describedby=' + tooltip.id + ']').length === 0) {
            tooltip.remove()
        }
    }
}

function shouldDisableDpoSaveButton(shouldDisable) {
    $('#' + modalDpoVars.dpoSaveID).prop('disabled', shouldDisable)
}

function isAnyDpoTooltips() {
    return $(modalDpoVars.tooltipsDpoClass).length !== 0
}

function setSelectDpoColumn(title, values) {
    let columnHeader = $('th:contains(' + title + ')')
    columnHeader.removeClass(modalDpoVars.dpoInputColumnClass)
    let index = columnHeader.index() + 1

    $(modalDpoVars.dpoTableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.prop(modalDpoVars.contentDpoEditableProp, false)

        let selector = $('<select>')

        for (let val of values) {
            const newOption = $('<option>', {
                value: val,
                text: val,
                selected: true
            });
            selector.append(newOption)
        }
        // currently a switch so if more values are added in the future, they can without converting if to switch
        switch (cell.text()) {
            case 'true':
                selector.find('option[value="true"]').prop('selected', true)
        }

        cell.html(selector)
    })
}

function revertSelectDpoColumn(title) {
    let index = $('th:contains(' + title + ')').index() + 1

    $(modalDpoVars.dpoTableID + ' td:nth-child(' + index + ')').each((i, c) => {
        let cell = $(c)

        cell.text(cell.find('option:selected').val())
    })
}

function emptyToolTipAndCloseDpoModal() {
    $(modalDpoVars.dpoTableID + ' td[aria-describedby*="tooltip"]').hide();
    $("#dpo-config-edit-modal").modal('hide');
    dpoModalCheck = 0

}

function checkDefaultValidation() {

    let checkRows = $(modalDpoVars.dpoTableID + ' tbody tr');
    let checkTrue = false;
    let count = 0;
    checkRows.each((i) => {
        let id = "#country-default-row-" + i;
        if ($(id)[0].value === "true") {
            checkTrue = true
            count = count + 1;
        }
    });
    console.log(checkTrue, count);
    return [checkTrue, count]
}