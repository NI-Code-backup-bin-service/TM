$(document).ready(function () {
    let token = $('input[name=csrfmiddlewaretoken]')[0].value
    let type = $('input[name=type]')[0].value;
    if (type === 'manage') {
        closeAddServiceModal()
        buildServicesTable(token)
    } else {
        closeAddGroupModal()
        buildServiceGroupTable(token)
    }
})

// SERVICE GROUP LOGIC:

function buildServiceGroupTable (token) {
    let table = $('#service-group-tbl').DataTable({
        ajax: {
            'url': 'paymentServicesSearch',
            'type': 'POST',
            'data': (d) => {
                return $.extend({}, d, {
                    csrfmiddlewaretoken: token,
                    'requestType': 'group'
                })
            }
        },
        deferRender: true,
        'autoWidth': false,
        'processing': true,
        "deferLoading": 0,
        'serverSide': true,
        'pageLength': 15,
        'lengthMenu': [
            [10, 15, 25, 50, 100, 500, -1],
            [10, 15, 25, 50, 100, 500, 'All']
        ],
        'columns': [
            { data: 'Name' },
            { data: 'ServiceCount', 'orderable': false,},
            {
                data: null,
                'orderable': false,
                'render': (cellData, cellType, rowData, _) => {
                    const name = rowData['Name'].replaceAll(' ', '-')
                    return `<div class="btn-group" role="group"> <button id="edit-group-${name}" class="btn btn-secondary btn-sm button-margin" data-button="edit-group">Edit</button>
                            <button id="delete-group-${name}" class="btn btn-outline-danger btn-sm button-margin" data-button="delete-group">Delete</button> </div>`
                }
            }
        ],
        'columnDefs': [{
            'targets': '_all',
            'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
                $(td).attr('id', `group-cell-${colIndex}-${rowData['Name'].replaceAll(' ', '-')}`)
            }
        }],
        'createdRow': (row, data, _) => {
            $(row).attr('id', `group-row-${data['Name'].replaceAll(' ', '-')}`)
        },

        initComplete: function() {
            const input = $('#service-group-tbl_filter input').attr('id', 'group-searchTerm').unbind();
            const self = this.api();
            const $searchButton = $('<button>')
                .text('Search')
                .attr('id', 'groupSearchButton')
                .attr('class', 'btn btn-primary button-margin')
                .click(function() {
                    self.search(input.val()).draw();
                });
            const $addButton = $('<button>')
                .text('Add Group')
                .attr('id', 'groupAddButton')
                .attr('class', 'btn btn-secondary button-margin')
                .click(function() {
                    showAddGroupModal()
                });
            $('#service-group-tbl_filter').append($searchButton).append($addButton);

            $('#' + 'group-searchTerm').keyup(function(event){
                if(event.key !== "Enter") return;
                event.preventDefault();
                $('#' + 'groupSearchButton').click();
            });
        }
    });

    $('#service-group-tbl tbody').on('click', 'button', function () {
        let row = table.row($(this).parents('tr'))
        let data = row.data()

        if ($(this).attr('data-button') === 'edit-group') {
            location.href = `/paymentServicesManagement?groupId=${data['Id']}`
            return
        }

        confirmDialog('Delete Group?', 'Are you sure you want to delete this group?', () => {
            $.ajax({
                url: 'paymentServicesDeleteGroup',
                data: {
                    groupId: data['Id'],
                    csrfmiddlewaretoken: token
                },
                method: 'POST',
                success: function () {
                    row.remove().draw();
                },
                error: function () {
                    displayWarningMessage("an error occurred deleting the group")
                }
            })
        })
    })
}

function closeAddGroupModal () {
    $('#' + 'add-group-modal').hide()
}

function showAddGroupModal () {
    $('#' + 'add-group-modal').show()
}

function addPaymentServiceGroup () {
    const name = $('#group-name-input').val()
    console.log('name: ' + name)
    if (!name || name.length === 0) {
        closeAddGroupModal()
        displayWarningMessage('error: cannot add a group with no name')
        return
    }

    $.ajax({
        url: 'paymentServicesAddGroup',
        data: {
            name: name,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: 'POST',
        success: function () {
            closeAddGroupModal()
        },
        error: function (data) {
            closeAddGroupModal()
            displayWarningMessage(`error: ${data.responseText}`)
        },
        complete(){
            $('#group-name-input').val('');
        }
    })
}

// INDIVIDUAL SERVICES LOGIC:

function buildServicesTable (token) {
    const url = new URL(window.location);
    const groupId = `${url.searchParams.get("groupId")}`

    let table = $('#service-tbl').DataTable({
        ajax: {
            'url': 'paymentServicesSearch',
            'type': 'POST',
            'data': (d) => {
                return $.extend({}, d, {
                    csrfmiddlewaretoken: token,
                    'requestType': 'service',
                    'groupId': groupId
                })
            }
        },
        deferRender: true,
        'autoWidth': false,
        'processing': true,
        "deferLoading": 0,
        'serverSide': true,
        'pageLength': 15,
        'lengthMenu': [
            [10, 15, 25, 50, 100, 500, -1],
            [10, 15, 25, 50, 100, 500, 'All']
        ],
        'columns': [
            { data: 'Name' },
            {
                data: null,
                'orderable': false,
                'render': (cellData, cellType, rowData, _) => {
                    const name = rowData['Name'].replaceAll(' ', '-')
                    return `<button id="delete-service-${name}" class="btn btn-outline-danger btn-sm button-margin" data-button="delete-service">Delete</button>`
                }
            }
        ],
        'columnDefs': [{
            'targets': '_all',
            'createdCell': (td, cellData, rowData, rowIndex, colIndex) => {
                $(td).attr('id', `service-cell-${colIndex}-${rowData['Name'].replaceAll(' ', '-')}`)
            }
        }],
        'createdRow': (row, data, _) => {
            $(row).attr('id', `service-row-${data['Name'].replaceAll(' ', '-')}`)
        },

        initComplete: function() {
            const input = $('#service-tbl_filter input').attr('id', 'service-searchTerm').unbind();
            const self = this.api();
            const $searchButton = $('<button>')
                .text('Search')
                .attr('id', 'serviceSearchButton')
                .attr('class', 'btn btn-primary button-margin')
                .click(function() {
                    self.search(input.val()).draw();
                });
            const $addButton = $('<button>')
                .text('Add Service')
                .attr('id', 'serviceAddButton')
                .attr('class', 'btn btn-secondary button-margin')
                .click(function() {
                    showAddServiceModal()
                });
            $('#service-tbl_filter').append($searchButton).append($addButton);

            $('#' + 'service-searchTerm').keyup(function(event){
                if(event.key !== "Enter") return;
                event.preventDefault();
                $('#' + 'serviceSearchButton').click();
            });
        }
    });

    $('#service-tbl tbody').on('click', 'button', function () {
        let row = table.row($(this).parents('tr'))
        let data = row.data()

        if ($(this).attr('data-button') === 'delete-service') {
            confirmDialog('Delete Service?', 'Are you sure you want to delete this service?', () => {
                $.ajax({
                    url: 'paymentServicesDelete',
                    data: {
                        serviceId: data['Id'],
                        csrfmiddlewaretoken: token
                    },
                    method: 'POST',
                    success: function () {
                        row.remove().draw()
                        hideWarning();
                    },
                    error: function () {
                        displayWarningMessage("an error occurred deleting the service")
                    }
                })
            })
        }
    })
}

function closeAddServiceModal () {
    $('#' + 'add-service-modal').hide()
}

function showAddServiceModal () {
    $('#' + 'add-service-modal').show()
}

function addPaymentService () {
    const url = new URL(window.location);
    const name = $('#service-name-input').val()
    if (!name || name.length === 0) {
        closeAddServiceModal()
        displayWarningMessage('error: cannot add a service with no name')
        return
    }

    $.ajax({
        url: 'paymentServicesAddService',
        data: {
            groupId: `${url.searchParams.get("groupId")}`,
            name: name,
            csrfmiddlewaretoken: $('input[name=csrfmiddlewaretoken]')[0].value
        },
        method: 'POST',
        success: function () {
            closeAddServiceModal()
        },
        error: function (data) {
            closeAddServiceModal()
            displayWarningMessage(`error: ${data.responseText}`)
        },
        complete() {
            $('#service-name-input').val('');
        }
    })
}
