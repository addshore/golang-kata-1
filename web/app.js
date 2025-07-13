// Simple client-side app for demo purposes only
// In a real app, you would use a backend API for CSV I/O

let books = [];
let magazines = [];
let authors = {};

async function loadCSV(path) {
  const res = await fetch(path);
  const text = await res.text();
  return text.split('\n').filter(Boolean);
}

function parseAuthors(lines) {
  const map = {};
  for (let i = 1; i < lines.length; i++) {
    const [email, firstName, lastName] = lines[i].split(';');
    if (email) map[email] = { email, firstName, lastName };
  }
  return map;
}

function parseBooks(lines) {
  const arr = [];
  for (let i = 1; i < lines.length; i++) {
    const [title, isbn, authors, description] = lines[i].split(';');
    if (isbn) arr.push({ title, isbn, authors: authors.split(','), description });
  }
  return arr;
}

function parseMagazines(lines) {
  const arr = [];
  for (let i = 1; i < lines.length; i++) {
    const [title, isbn, authors, publishedAt] = lines[i].split(';');
    if (isbn) arr.push({ title, isbn, authors: authors.split(','), publishedAt });
  }
  return arr;
}

async function loadAll() {
  const [a, b, m] = await Promise.all([
    loadCSV('/api/csv/authors.csv'),
    loadCSV('/api/csv/books.csv'),
    loadCSV('/api/csv/magazines.csv'),
  ]);
  authors = parseAuthors(a);
  books = parseBooks(b);
  magazines = parseMagazines(m);
  showAll();
}

function showAll() {
  let html = '<h2>Books</h2>';
  for (const b of books) html += renderBook(b);
  html += '<h2>Magazines</h2>';
  for (const m of magazines) html += renderMagazine(m);
  document.getElementById('results').innerHTML = html;
}

function renderBook(b) {
  return `<div><b>${b.title}</b> (ISBN: ${b.isbn})<br>Authors: ${b.authors.map(a=>resolveAuthor(a)).join(', ')}<br>${b.description}</div>`;
}
function renderMagazine(m) {
  return `<div><b>${m.title}</b> (ISBN: ${m.isbn})<br>Authors: ${m.authors.map(a=>resolveAuthor(a)).join(', ')}<br>Published: ${m.publishedAt}</div>`;
}
function resolveAuthor(email) {
  const a = authors[email];
  return a ? `${a.firstName} ${a.lastName}` : email;
}

function search() {
  const q = document.getElementById('searchInput').value.trim();
  if (!q) return;
  let found = books.filter(b => b.isbn === q || b.authors.includes(q));
  found = found.concat(magazines.filter(m => m.isbn === q || m.authors.includes(q)));
  let html = `<h2>Search Results for "${q}"</h2>`;
  if (found.length === 0) html += '<div>No results.</div>';
  for (const item of found) {
    html += item.description ? renderBook(item) : renderMagazine(item);
  }
  document.getElementById('results').innerHTML = html;
}

function sortByTitle() {
  const all = books.concat(magazines);
  all.sort((a, b) => a.title.localeCompare(b.title));
  let html = '<h2>All Books & Magazines Sorted by Title</h2>';
  for (const item of all) {
    html += item.description ? renderBook(item) : renderMagazine(item);
  }
  document.getElementById('results').innerHTML = html;
}

function showAddForm() {
  document.getElementById('formDiv').style.display = '';
  document.getElementById('formDiv').innerHTML = `
    <h3>Add to Library</h3>
    <select id="addType">
      <option value="book">Book</option>
      <option value="magazine">Magazine</option>
      <option value="author">Author</option>
    </select>
    <div id="addFields"></div>
    <button onclick="submitAdd()">Add</button>
    <button onclick="hideAddForm()">Cancel</button>
    <div id="addMsg"></div>
  `;
  document.getElementById('addType').onchange = updateAddFields;
  updateAddFields();
}
function hideAddForm() {
  document.getElementById('formDiv').style.display = 'none';
}
function updateAddFields() {
  const t = document.getElementById('addType').value;
  let html = '';
  if (t === 'book') {
    html += 'Title: <input id="addTitle"><br>';
    html += 'ISBN: <input id="addISBN"><br>';
    html += 'Authors (comma emails): <input id="addAuthors"><br>';
    html += 'Description: <input id="addDesc"><br>';
  } else if (t === 'magazine') {
    html += 'Title: <input id="addTitle"><br>';
    html += 'ISBN: <input id="addISBN"><br>';
    html += 'Authors (comma emails): <input id="addAuthors"><br>';
    html += 'Published At: <input id="addPub"><br>';
  } else {
    html += 'Email: <input id="addEmail"><br>';
    html += 'First Name: <input id="addFirst"><br>';
    html += 'Last Name: <input id="addLast"><br>';
  }
  document.getElementById('addFields').innerHTML = html;
}
async function submitAdd() {
  const t = document.getElementById('addType').value;
  let msg = '';
  let payload = { type: t, data: {} };
  if (t === 'book') {
    payload.data = {
      title: document.getElementById('addTitle').value,
      isbn: document.getElementById('addISBN').value,
      authors: document.getElementById('addAuthors').value,
      description: document.getElementById('addDesc').value
    };
  } else if (t === 'magazine') {
    payload.data = {
      title: document.getElementById('addTitle').value,
      isbn: document.getElementById('addISBN').value,
      authors: document.getElementById('addAuthors').value,
      publishedAt: document.getElementById('addPub').value
    };
  } else {
    payload.data = {
      email: document.getElementById('addEmail').value,
      firstName: document.getElementById('addFirst').value,
      lastName: document.getElementById('addLast').value
    };
  }
  const resp = await fetch('/api/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  });
  if (resp.ok) {
    msg = 'Saved!';
    await loadAll();
  } else {
    msg = 'Error saving.';
  }
  document.getElementById('addMsg').innerText = msg;
}
