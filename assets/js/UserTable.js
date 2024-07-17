// noinspection ES6ModulesDependencies
const ce = React.createElement;
let addedRowId = -1;

function generateNewRandomPin(users) {
    let newPin = Math.floor(Math.random() * (99999));
    let pin = "" + newPin;
    pin = pin.padStart(5, "0");

    if (users.findIndex(function (val, index, obj) {
        return val.PIN == newPin;
    }) > -1) {
        pin = generateNewRandomPin(users)
    }
    return pin
}


class UserTable extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            Users: props.Users,
            Modules: props.Modules,
            FriendlyModules: props.FriendlyModules,
            DeletedUsers: new Set(),
            ModifiedUsers: new Set(),
            SaveUsers: props.UserSave,
            Page: 0,
            PerPage: props.PerPage,
            hasSavePermission: props.hasSavePermission,
            hasPasswordPermission: props.hasPasswordPermission
        };
        this.save = this.save.bind(this);
        this.addModified = this.addModified.bind(this);
        this.addDeleted = this.addDeleted.bind(this);
        this.addDeletedNew = this.addDeletedNew.bind(this);
        this.nextPage = this.nextPage.bind(this);
        this.previousPage = this.previousPage.bind(this);
        this.addRow = this.addRow.bind(this);
        this.deleteRow = this.deleteRow.bind(this);
        this.deleteRowNew = this.deleteRowNew.bind(this);
        this.genNewRandomPin = this.genNewRandomPin.bind(this);
        this.selectAll = this.selectAll.bind(this);
        this.checkPinExists = this.checkPinExists.bind(this);
        this.checkPins = this.checkPins.bind(this);
    }

    genNewRandomPin() {
        return generateNewRandomPin(this.state.Users)
    }

    checkPins() {
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (val.PIN.length < 4) {
                return true
            }

        });

        return index <= -1;
    }

    checkDuplicates() {
        const pins = new Set();
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (pins.has(val.PIN)) {
                return true
            } else {
                pins.add(val.PIN)
            }
        });
        return index > -1;
    }

    checkDuplicateName() {
        const names = new Set();
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (Array.from(names).some(item => item.toLowerCase() === val.Username.toLowerCase())) {
                return true
            } else {
                names.add(val.Username)
            }
        });
        return index > -1;
    }

    checkNames() {
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (val.Username.length == 0 && val.UserId > 0) {
                return true
            }
        });

        return index <= -1;
    }



    checkSuperPins() {
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (superPins.includes(val.PIN)) {
                return true
            }

        });

        return index <= -1;
    }

    save() {
        this.setState({ hasSavePermission: false });
        if (!this.checkPins()) {
            this.props.DisplayWarningMessage("PIN(s) must be 4 or 5 digits.")
        } else if (!this.checkSuperPins()) {
            this.props.DisplayWarningMessage("PIN(s) must not match supervisor PIN(s)")
        } else if (this.checkDuplicates()) {
            this.props.DisplayWarningMessage("PIN(s) must be unique")
        } else if (!this.checkNames()) {
            this.props.DisplayWarningMessage("Username(s) name cannot be empty")
        } else if (this.checkDuplicateName()) {
            this.props.DisplayWarningMessage("Duplicate user")
        } else {
            let success = this.state.SaveUsers(this.state.Users, this.state.DeletedUsers, this.state.ModifiedUsers)
            if (success) {
                this.state.ModifiedUsers = [];
            }
        }
        this.setState({ hasSavePermission: true });
    }

    nextPage() {
        this.setState({ Page: this.state.Page + 1 })
    }

    previousPage() {
        this.setState({ Page: this.state.Page - 1 })
    }

    addModified(id) {
        const modified = this.state.ModifiedUsers;
        if (!this.state.DeletedUsers.has(id)) {
            modified.add(parseInt(id, 10));
            this.setState({ ModifiedUsers: modified })
        }
    }

    addDeleted(id) {
        const deleted = this.state.DeletedUsers;
        if (this.state.ModifiedUsers.has(id)) {
            const modified = this.state.ModifiedUsers;
            modified.delete(parseInt(id, 10))
            this.setState({ ModifiedUsers: modified })
        }
        deleted.add(parseInt(id, 10))
        return deleted
    }

    addDeletedNew(user) {
        const deleted = this.state.DeletedUsers;
        if (this.state.ModifiedUsers.has(user.UserId)) {
            const modified = this.state.ModifiedUsers;
            modified.delete(parseInt(user.UserId, 10))
            this.setState({ ModifiedUsers: modified })
        }
        deleted.add(user)
        return deleted
    }

    addRow() {
        const users = this.state.Users;

        users.push({
            UserId: addedRowId--,
            Username: '',
            PIN: this.genNewRandomPin(),
            Modules: ['sale', 'gratuitySale', 'X-Read', 'Z-Read']
        });
        this.setState({
            Users: users
        });
        if (users.length > (this.state.Page + 1) * this.state.PerPage) {
            this.nextPage()
        }
    }

    deleteRow(userId) {
        const users = this.state.Users.filter(
            function (item) {
                return item.UserId != userId;
            });

        if (userId > -1) {
            this.state.Users.forEach(user => {
                if (user.UserId == userId) {
                    this.setState({
                        Users: users,
                        DeletedUsers: this.addDeleted(userId)
                    });

                }
            });
        } else {
            this.setState({
                Users: users
            });
        }
    }

    deleteRowNew(userId) {
        const users = this.state.Users.filter(
            function (item) {
                return item.UserId != userId;
            });

        if (userId > -1) {
            this.state.Users.forEach(user => {
                if (user.UserId == userId) {
                    delete user.SiteId
                    delete user.TidId
                    this.setState({
                        Users: users,
                        DeletedUsers: this.addDeletedNew(user)
                    });

                }
            });
        } else {
            this.setState({
                Users: users
            });
        }
    }

    selectAll(id) {
        const users = this.state.Users;
        let index = users.findIndex(function (val, index, obj) {
            if (val.UserId == id) {
                return true
            }

        });
        if (index > -1) {
            users[index].Modules = this.state.Modules;
            this.setState({ Users: users });
            this.addModified(id)
        }
    }

    checkPinExists(newPin, userID) {
        if (this.state.Users.findIndex(function (val, index, obj) {
            return val.PIN == newPin && val.UserId != userID
        }) > -1) {
            return true
        }
        return false
    }

    render() {
        const header = UserTableHeader({
            key: "header",
            Modules: this.state.Modules,
            FriendlyModules: this.state.FriendlyModules
        });

        const body = ce(UserTableBody, {
            key: "body",
            Users: this.state.Users,
            Modules: this.state.Modules,
            AddModified: this.addModified,
            Page: this.state.Page,
            NextPage: this.nextPage,
            PreviousPage: this.previousPage,
            PerPage: this.props.PerPage,
            AddRow: this.addRow,
            DeleteRow: this.deleteRowNew,
            SelectAll: this.selectAll,
            CheckPinExists: this.checkPinExists,
            hasSavePermission: this.state.hasSavePermission,
            hasPasswordPermission: this.state.hasPasswordPermission
        }, null);

        const elements = [header, body];

        const table = ce('table', { key: "table", className: 'table table-striped table-sm tms-user-table' }, elements);
        const tableDiv = ce('div', { key: "tableContainer", className: "tms-user-table-container" }, table);

        const saveButton = ce('button', {
            id: "user-management-save",
            key: "save-button",
            className: 'btn btn-primary spaced-bottom',
            onClick: this.save,
            disabled: (!this.state.hasSavePermission || !this.state.hasPasswordPermission)
        }, "Save");

        let prevDisabled = false;
        if (this.state.Page <= 0) {
            prevDisabled = true;
        }

        let nextDisabled = false;
        if ((this.state.Page + 1) * this.props.PerPage >= this.state.Users.length) {
            nextDisabled = true;
        }

        const prevPage = ce('button', {
            id: "user-management-previous",
            key: "prev-button",
            className: 'btn btn-primary spaced-bottom float-start',
            onClick: this.previousPage,
            disabled: prevDisabled
        }, "Previous");

        const pageLabel = ce('label', {
            className: 'tms-user-page-label',
            key: "tms-user-page-label"
        },
            "Page " + (this.state.Page + 1) + " of " +
            (this.state.Users.length / this.state.PerPage < 1 ? 1 : Math.floor(this.state.Users.length / this.props.PerPage) + (this.state.Users.length % this.state.PerPage != 0 ? 1 : 0)))

        const nextPage = ce('button', {
            id: "user-management-next",
            key: "next-button",
            className: 'btn btn-primary spaced-bottom float-end',
            onClick: this.nextPage,
            disabled: nextDisabled
        }, "Next");

        return ce('div', { key: "userTableDiv" }, [prevPage, pageLabel, nextPage, tableDiv, saveButton])
    }
}

function UserTableHeader(props) {
    const columns = [
        ce('th', { key: "Username", style: { width: "140px" }, className: 'tms-user-module-cell' }, 'Username'),
        ce('th', { key: "PIN", style: { width: "130px" }, className: 'tms-user-module-cell' }, 'PIN')
    ];

    for (let i = 0; i < props.Modules.length; i++) {
        columns.push(ce('th', {
            key: props.Modules[i],
            className: "tms-user-module-cell",
            style: { width: "90px" }
        }, props.FriendlyModules[i])
        );
    }

    columns.push(ce('th', { key: "select-all-column", style: { width: "90px" } }, null));
    columns.push(ce('th', { key: "button-column", style: { width: "90px" } }, null));

    return ce('thead', { key: "thead" }, [ce('tr', { key: "header-row" }, columns)])
}

class UserTableBody extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            Modules: props.Modules,
            AddDeleted: props.AddDeleted,
            AddDeletedNew: props.AddDeletedNew,
            AddModified: props.AddModified,
        };
    }

    render() {
        let rows = [];
        let id = 0;
        let maxIndex = this.props.Users.length;
        if (this.props.PerPage * (this.props.Page + 1) < this.props.Users.length) {
            maxIndex = this.props.PerPage * (this.props.Page + 1)
        }
        for (let i = this.props.Page * this.props.PerPage; i < maxIndex; i++) {
            rows.push(ce(UserTableRow, {
                HtmlId: ++id,
                key: this.props.Users[i].UserId,
                User: this.props.Users[i],
                Modules: this.state.Modules,
                IsLastRow: i === this.props.Users.length - 1,
                AddRow: this.props.AddRow,
                DeleteRow: this.props.DeleteRow,
                AddModified: this.state.AddModified,
                SelectAll: this.props.SelectAll,
                CheckPinExists: this.props.CheckPinExists,
                hasSavePermission: this.props.hasSavePermission,
                hasPasswordPermission: this.props.hasPasswordPermission
            }, null));
        }

        return ce("tbody", { key: "tbody" }, rows)
    }

}

class UserTableRow extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            HtmlId: props.HtmlId,
            User: props.User,
            Modules: props.Modules,
            AddRow: props.AddRow,
            DeleteRow: props.DeleteRow,
            IsLastRow: props.IsLastRow,
            ModifiedRow: props.AddModified,
            SelectAll: props.SelectAll,
            CheckPinExists: props.CheckPinExists,
            hasSavePermission: props.hasSavePermission,
            hasPasswordPermission: props.hasPasswordPermission,
            PinExists: false,
            ShowPIN: false
        };
        this.modulesChanged = this.modulesChanged.bind(this);
        this.userNameChanged = this.userNameChanged.bind(this);
        this.userPINChanged = this.userPINChanged.bind(this);
        this.rowButtonClick = this.rowButtonClick.bind(this);
        this.selectAll = this.selectAll.bind(this);
        this.togglePIN = this.togglePIN.bind(this)

    }

    userNameChanged(e) {
        const user = this.state.User;
        user.Username = e.target.value;
        this.setState({ User: user });
        this.state.ModifiedRow(this.state.User.UserId)
    }

    userPINChanged(e) {
        const regex = /^[0-9\b]+$/;
        if (e.target.value !== "" && !regex.test(e.target.value)) {
            e.preventDefault();
            return
        }
        const user = this.state.User;
        user.PIN = e.target.value;
        this.setState({ User: user, PinExists: this.state.CheckPinExists(e.target.value, this.state.User.UserId) });
        this.state.ModifiedRow(this.state.User.UserId)
    }

    togglePIN(e) {
        this.setState({ ShowPIN: !this.state.ShowPIN })
    }

    modulesChanged(e) {
        const module = e.target.name;
        const checked = e.target.checked;
        const user = this.state.User;
        if (checked) {
            user.Modules.push(module)
        } else {
            user.Modules = this.state.User.Modules.filter(
                function (item) {
                    return item !== module
                }
            )
        }
        this.setState({ User: user });
        this.state.ModifiedRow(this.state.User.UserId)
    }

    rowButtonClick(e) {
        if (this.state.IsLastRow) {
            this.setState({
                IsLastRow: false
            });
            this.state.AddRow()

        } else {
            this.state.DeleteRow(e.target.name)
        }
    }


    selectAll(e) {
        this.state.SelectAll(e.target.name)
    }

    render() {
        const id = this.state.HtmlId;

        //console.log("hasSavePermission="+this.state.hasSavePermission+ " , hasPasswordPermission=" + this.state.hasPasswordPermission)
        //console.log(this.state);

        var isToggleDisabled = (!this.state.hasSavePermission || !this.state.hasPasswordPermission);

        const togglePIN = ce('img', {
            key: 'togglePIN',
            onClick: isToggleDisabled ? null : this.togglePIN,
            name: this.state.User.UserId,
            id: "Row" + id + "-showPwdBtn",
            src: "/assets/images/show-password-icon.png",
            className: "left-spaced togglePIN",
            style: isToggleDisabled ? { display: 'inline-block', opacity: 0.5 } : { display: 'inline-block' },
            disabled: isToggleDisabled
        });

        const cells = [
            ce('td', { key: "Username" }, ce('input', {
                id: "user-management-username-" + id,
                key: "Username",
                type: 'text',
                value: this.state.User.Username,
                size: 15,
                onChange: this.userNameChanged,
                className: "react-table-input",
                maxLength: 10
            }, null)),
            ce('td', { key: "PIN" }, [ce('input', {
                id: "user-management-pin-" + id,
                key: "PIN",
                type: this.state.ShowPIN ? "text" : "password",
                size: 5,
                value: this.state.User.PIN,
                onChange: this.userPINChanged,
                maxLength: 5,
                pattern: "\\d",
                className: "react-table-input" + (this.state.PinExists ? ' error' : '')
            }, null), togglePIN])
        ];
        for (let i = 0; i < this.state.Modules.length; i++) {
            if (this.state.User.Modules.includes(this.state.Modules[i])) {
                var box = ce('input', {
                    id: "user-management-" + this.state.Modules[i] + "-" + id,
                    key: this.state.Modules[i],
                    type: 'checkbox',
                    name: this.state.Modules[i],
                    checked: true,
                    onChange: this.modulesChanged,
                    className: "form-checkbox profile-checkbox"
                }, null);
                cells.push(ce('td', { key: this.state.Modules[i], className: 'tms-user-module-cell' }, box))
            } else {
                var box = ce('input', {
                    id: "user-management-" + this.state.Modules[i] + "-" + id,
                    key: this.state.Modules[i],
                    type: 'checkbox',
                    name: this.state.Modules[i],
                    checked: false,
                    onChange: this.modulesChanged,
                    className: "form-checkbox profile-checkbox"
                }, null);
                cells.push(ce('td', { key: this.state.Modules[i], className: 'tms-user-module-cell' }, box))
            }
        }


        const selectAll = ce('button', {
            id: "user-management-select-all-" + id,
            key: "select-all",
            onClick: this.selectAll,
            name: this.state.User.UserId,
            className: "btn btn-sm btn-secondary"
        }, 'Select All');
        cells.push(ce('td', { key: "select-all-cell" }, selectAll));

        const button = ce('button', {
            id: "user-management-add-delete-" + id,
            key: "add/delete",
            onClick: this.rowButtonClick,
            name: this.state.User.UserId,
            className: "btn btn-sm " + (this.state.IsLastRow ? "btn-secondary" : "btn-danger")
        }, this.state.IsLastRow ? 'Add' : 'Delete');



        cells.push(ce('td', { key: "button-cell" }, button));

        return ce('tr', { key: this.state.User.UserId }, cells)
    }
}

//# sourceURL=reactUserTable.js