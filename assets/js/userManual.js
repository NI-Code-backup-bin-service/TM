function printDocument() {
    const tab = $("#tmsTabs").find(".active")[0].id.replace('#','')
    if ($('#'+tab).children('div.tab-content').length >= 1){
        const nestedTab = $('#'+tab).children('div.tab-content').find('.active').find('.active')[0].id.replace('#','')
        if ( $('#' + nestedTab).children('embed')[0] == undefined){
            console.log($('#'+tab).children('div.tab-content').find('.active').find('.active')[0].id.replace('#',''))
            printJS(nestedTab, 'html')
        } else {
            printJS($('#' + nestedTab).children('embed')[0].src)
        }
    } else if ($('#'+tab).children('embed')[0] == undefined){
        printJS(tab, 'html')
    }else{
        printJS($('#'+tab).children('embed')[0].src)
    }
}
async function exportPDF() {
    const { PDFDocument } = PDFLib
    const pdfDoc = await PDFDocument.create();
    let urls = []
    $('#tmsContent').children().each(function (ind, element) {
        console.log(element);
        $(element).children('embed').each(function (ind, element){
            urls.push(element.src)
            console.log(urls.length)
        });
    });
    for (let url of urls){
        let pdfBytes = await fetch(url).then((res) => res.arrayBuffer());
        let pdf = await PDFDocument.load(pdfBytes);
        let pagesArray = await pdfDoc.copyPages(pdf, pdf.getPageIndices())
        console.log("AddPages")
        for await (let page of pagesArray) {
            pdfDoc.addPage(page);
        }
    }
    // Serialize the PDFDocument to bytes (a Uint8Array)
    const pdfSaveBytes = await pdfDoc.save()
    // Trigger the browser to download the PDF document
    download(pdfSaveBytes, "tmsUserManual.pdf", "application/pdf");
}