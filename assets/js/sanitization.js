/*
    These js functions are to cleanse html and help prevent XSS attacks.
    To use this file, you must first include "DOMPurify/purify.min.js"
*/

function sanitizeHTML(inHtml) {
    const config = {
        ADD_ATTR: [
            'onclick',
            'onerror',
            'onchange',
            'ondrop',
            'ondragover',
            'ondragstart',
            'draggable',
            'rv-value',
            'rv-each-transaction',
            'rv-each-limittype',
            'rv-i',
            'rv-tid',
            'rv-on-click',
            'rv-each-limit'
        ],
        FORCE_BODY: true,
        ADD_TAGS: ['script'],
        SAFE_FOR_JQUERY: true,
    }

    return DOMPurify.sanitize(inHtml, config);
}

/*
    If the code being sanitised contains table elements (<tr>, <td>, etc) but isn't enclosed in <table> tags these will be flagged
    as invalid HTML and removed. This wraps the html in table tags and then removes them post-sanitization to avoid this.
 */
function sanitizeTableHTML(inHtml) {
    inHtml = '<table><tbody>' + inHtml + '</tbody></table>'
    const config = {
        ADD_ATTR: [
            'onclick',
            'onerror',
            'onchange',
            'ondrop',
            'ondragover',
            'ondragstart',
            'draggable'
        ],
        SAFE_FOR_JQUERY: true,
    }

    let clean = DOMPurify.sanitize(inHtml, config);

    clean = clean.substring(14);
    clean = clean.substring(0, clean.length-17);
    return clean;
}