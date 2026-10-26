// Run from the browser skill's Node session with a connected Tab handle.
// Reads the rendered official page, never hidden state or credential stores.
export async function captureReference(tab, sourceUrl) {
  if (await tab.url() !== sourceUrl) await tab.goto(sourceUrl);
  await tab.playwright.getByRole('heading', { level: 1 }).first().waitFor({state:'visible', timeoutMs:15000});
  return await tab.playwright.evaluate((url) => {
    const root = document.querySelector('.layout-right');
    const title = root?.querySelector('h1')?.innerText;
    if (!root || !title) throw new Error('Official document body is unavailable');
    const normalize = s => s.replace(/\s+/g, ' ').trim();
    let section = '', location = '', sectionUrl = url, authenticationBranch = '', example = null;
    const tables = [], notes = [], examples = [], links = [], seen = new Set();
    for (const element of root.querySelectorAll('h2,h3,h4,table,p,li,pre')) {
      if (/^H[234]$/.test(element.tagName)) {
        if (element.tagName === 'H2') { section = normalize(element.innerText); location = ''; authenticationBranch = ''; example = null; }
        else if (['Header','Query','Body','Path'].includes(element.innerText.trim())) location = element.innerText.trim().toLowerCase();
        sectionUrl = element.id ? `${url}#${element.id}` : sectionUrl;
        if (/请求.*[例样]|响应.*[例样]/.test(section) && element.tagName==='H2') {
          // A heading's parent may contain only that heading, while examples
          // live in sibling application-kind subsections. Accumulate visible
          // code blocks until the next H2 instead of guessing a container.
          example = {section, sourceUrl:sectionUrl, text:element.innerText.trim()};
          examples.push(example);
        } else if (example) {
          example.text += '\n' + element.innerText.trim();
        }
      } else if (element.tagName === 'PRE' && example) {
        example.text += '\n' + element.innerText.trim();
      } else if (element.tagName === 'TABLE') {
        const rows = [], spans = [];
        for (const tr of element.querySelectorAll('tr')) {
          const row = []; let column = 0;
          const copySpans = () => { while (spans[column]?.left > 0) { row[column] = spans[column].text; spans[column].left--; column++; } };
          for (const cell of tr.querySelectorAll(':scope > td,:scope > th')) {
            copySpans();
            const text = normalize(cell.innerText);
            for (let j = 0; j < (cell.colSpan || 1); j++) {
              row[column] = text;
              if (cell.rowSpan > 1) spans[column] = {text, left:cell.rowSpan - 1};
              column++;
            }
          }
          copySpans(); rows.push(row);
        }
        tables.push({section, location, sourceUrl:sectionUrl, ...(location === 'header' && authenticationBranch ? {authenticationBranch} : {}), rows});
      } else if (!element.closest('table') && !/[例样]/.test(section)) {
        const text = normalize(element.innerText);
        if (/^(Service Account|API客户端|OAuth客户端)方式[:：]?$/.test(text)) authenticationBranch = text;
        if (text && !seen.has(text)) { notes.push({section:section || '文档说明', sourceUrl:sectionUrl, text}); seen.add(text); }
      }
    }
    for (const a of root.querySelectorAll('a[href]')) {
      if (a.href.includes('/doc/') && a.innerText.trim()) links.push({title:normalize(a.innerText),sourceUrl:a.href});
    }
    const updatedAt = root.innerText.match(/更新时间[:：]\s*([\d-]+\s+[\d:]+)/)?.[1];
    return {title,sourceUrl:url,updatedAt,checkedAt:new Date().toISOString().slice(0,10),renderedText:root.innerText,status:'browser-read',tables,notes,examples,links};
  }, sourceUrl);
}
