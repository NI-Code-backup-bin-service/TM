// noinspection ES6ModulesDependencies
const ce = React.createElement;


function userToCsv(headers, rows, filename) {
    let headerRow = headers.join(",") + "\n";
    let datarows = "";
    rows.forEach(function (value, index) {
        datarows += rows[index].join(",") + "\n";
    });

    var csvFile = new Blob([headerRow + datarows], {type: "text/csv"})
    var link = document.createElement('a');
    link.href = encodeURI(window.URL.createObjectURL(csvFile));
    link.download = filename + ".csv";
    link.click();
}

class ExportTable extends React.Component {
    constructor(props) {
        super(props);
        this.export = this.export.bind(this)
    }

    export(){
        userToCsv(this.props.Columns, this.props.Rows, this.props.Filename)
    }

    render() {
        let headers = [];
        for (let i = 0; i < this.props.Columns.length; i++) {
            headers.push(ce("th", {key: this.props.Columns[i]}, this.props.Columns[i]));
        }
        let head = ce('thead', {key: "Head"}, ce('tr', null, headers));

        let rows = [];
        for (let i = 0; i < this.props.Rows.length; i++) {
            let cells = [];
            for (let j = 0; j < this.props.Rows[i].length; j++) {
                cells.push(ce("td", {key: i + this.props.Rows[i][j]}, this.props.Rows[i][j]));
            }
            rows.push(ce("tr", {key: i}, cells))
        }
        let body = ce('tbody', {key: "Body"}, rows);

        let exportButton =
            ce('button', {
                id: 'export-button',
                key: "export-button",
                className: 'btn btn-primary spaced-bottom',
                onClick: this.export
            }, "Export");
        return ce('div',{className: "padded"},[ce('table', {key:"export-table", className: 'table table-striped table-sm'}, [head, body]), exportButton])
    }

}