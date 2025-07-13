package main

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Library Application</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #333;
            text-align: center;
        }
        .section {
            margin: 30px 0;
        }
        h2 {
            color: #555;
            border-bottom: 2px solid #ddd;
            padding-bottom: 10px;
        }
        .items {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
            gap: 20px;
            margin-top: 20px;
        }
        .item {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .item h3 {
            margin-top: 0;
            color: #333;
        }
        .isbn {
            color: #666;
            font-size: 0.9em;
            margin: 5px 0;
        }
        .authors {
            color: #888;
            font-size: 0.9em;
            margin: 10px 0;
        }
        .description {
            color: #444;
            line-height: 1.5;
            margin-top: 10px;
        }
        .published {
            color: #666;
            font-size: 0.9em;
            margin-top: 10px;
        }
        .author-name {
            font-weight: bold;
            color: #2c5aa0;
        }
        .stats {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
            margin-bottom: 30px;
        }
        .stats span {
            display: inline-block;
            margin: 0 20px;
            font-size: 1.2em;
        }
        .stats strong {
            color: #2c5aa0;
        }
        .search-box {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 30px;
            text-align: center;
        }
        .search-box form {
            display: inline-block;
        }
        .search-inputs {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .or-divider {
            color: #999;
            font-weight: bold;
            padding: 0 5px;
        }
        .search-box input[type="text"] {
            padding: 10px 15px;
            border: 1px solid #ddd;
            border-radius: 4px;
            width: 300px;
            font-size: 16px;
        }
        .search-box button {
            padding: 10px 20px;
            background-color: #2c5aa0;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 16px;
            cursor: pointer;
            margin-left: 10px;
        }
        .search-box button:hover {
            background-color: #1e3d6f;
        }
        .clear-search {
            display: inline-block;
            margin-left: 10px;
            color: #666;
            text-decoration: none;
        }
        .clear-search:hover {
            text-decoration: underline;
        }
        .search-results {
            margin-top: 15px;
            color: #666;
        }
        .view-toggle {
            text-align: center;
            margin-bottom: 30px;
        }
        .view-toggle span {
            margin-right: 10px;
            color: #666;
        }
        .view-link {
            display: inline-block;
            padding: 8px 16px;
            margin: 0 5px;
            background-color: #f0f0f0;
            color: #333;
            text-decoration: none;
            border-radius: 4px;
            transition: background-color 0.3s;
        }
        .view-link:hover {
            background-color: #e0e0e0;
        }
        .view-link.active {
            background-color: #2c5aa0;
            color: white;
        }
        .type-badge {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.8em;
            margin-left: 10px;
        }
        .type-book {
            background-color: #4CAF50;
            color: white;
        }
        .type-magazine {
            background-color: #FF9800;
            color: white;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Library Application</h1>
        
        <div class="stats">
            <span><strong>{{.BookCount}}</strong> Books</span>
            <span><strong>{{.MagazineCount}}</strong> Magazines</span>
            <span><strong>{{.AuthorCount}}</strong> Authors</span>
        </div>

        <div class="search-box">
            <form method="GET" action="/">
                <div class="search-inputs">
                    <input type="text" name="isbn" placeholder="Search by ISBN..." value="{{.SearchISBN}}" />
                    <span class="or-divider">OR</span>
                    <input type="text" name="author" placeholder="Search by author email..." value="{{.SearchAuthor}}" />
                    <input type="hidden" name="view" value="{{.ViewMode}}" />
                    <input type="hidden" name="sort" value="{{.SortDir}}" />
                    <button type="submit">Search</button>
                    {{if or .SearchISBN .SearchAuthor}}
                    <a href="/?view={{.ViewMode}}&sort={{.SortDir}}" class="clear-search">Clear Search</a>
                    {{end}}
                </div>
            </form>
            {{if .SearchISBN}}
            <div class="search-results">
                Showing results for ISBN: <strong>{{.SearchISBN}}</strong>
            </div>
            {{else if .SearchAuthor}}
            <div class="search-results">
                Showing results for author: <strong>{{.SearchAuthor}}</strong>
            </div>
            {{end}}
        </div>

        <div class="view-toggle">
            <span>View mode:</span>
            <a href="/?view=separate&sort={{.SortDir}}{{if .SearchISBN}}&isbn={{.SearchISBN}}{{end}}{{if .SearchAuthor}}&author={{.SearchAuthor}}{{end}}" 
               class="view-link {{if eq .ViewMode "separate"}}active{{end}}">Separate</a>
            <a href="/?view=combined&sort={{.SortDir}}{{if .SearchISBN}}&isbn={{.SearchISBN}}{{end}}{{if .SearchAuthor}}&author={{.SearchAuthor}}{{end}}" 
               class="view-link {{if eq .ViewMode "combined"}}active{{end}}">Combined (Sorted by Title)</a>
            
            {{if eq .ViewMode "combined"}}
            <span style="margin-left: 20px;">Sort:</span>
            <a href="/?view=combined&sort=asc{{if .SearchISBN}}&isbn={{.SearchISBN}}{{end}}{{if .SearchAuthor}}&author={{.SearchAuthor}}{{end}}" 
               class="view-link {{if eq .SortDir "asc"}}active{{end}}">A-Z</a>
            <a href="/?view=combined&sort=desc{{if .SearchISBN}}&isbn={{.SearchISBN}}{{end}}{{if .SearchAuthor}}&author={{.SearchAuthor}}{{end}}" 
               class="view-link {{if eq .SortDir "desc"}}active{{end}}">Z-A</a>
            {{end}}
        </div>

        {{if eq .ViewMode "combined"}}
        <div class="section">
            <h2>All Items (Sorted by Title {{if eq .SortDir "asc"}}A-Z{{else}}Z-A{{end}})</h2>
            <div class="items">
                {{range .CombinedItems}}
                <div class="item">
                    <h3>
                        {{.Title}}
                        {{if eq .Type "book"}}
                        <span class="type-badge type-book">Book</span>
                        {{else}}
                        <span class="type-badge type-magazine">Magazine</span>
                        {{end}}
                    </h3>
                    <div class="isbn">ISBN: {{.ISBN}}</div>
                    <div class="authors">
                        Authors: 
                        {{range $i, $email := .Authors}}
                            {{if $i}}, {{end}}
                            {{$author := index $.AuthorsMap $email}}
                            {{if $author}}
                                <span class="author-name">{{$author.FirstName}} {{$author.LastName}}</span>
                            {{else}}
                                {{$email}}
                            {{end}}
                        {{end}}
                    </div>
                    {{if eq .Type "book"}}
                    <div class="description">{{.Description}}</div>
                    {{else}}
                    <div class="published">Published: {{.PublishedAt.Format "02.01.2006"}}</div>
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
        {{else}}
        <div class="section">
            <h2>Books</h2>
            <div class="items">
                {{range .Books}}
                <div class="item">
                    <h3>{{.Title}}</h3>
                    <div class="isbn">ISBN: {{.ISBN}}</div>
                    <div class="authors">
                        Authors: 
                        {{range $i, $email := .Authors}}
                            {{if $i}}, {{end}}
                            {{$author := index $.AuthorsMap $email}}
                            {{if $author}}
                                <span class="author-name">{{$author.FirstName}} {{$author.LastName}}</span>
                            {{else}}
                                {{$email}}
                            {{end}}
                        {{end}}
                    </div>
                    <div class="description">{{.Description}}</div>
                </div>
                {{end}}
            </div>
        </div>

        <div class="section">
            <h2>Magazines</h2>
            <div class="items">
                {{range .Magazines}}
                <div class="item">
                    <h3>{{.Title}}</h3>
                    <div class="isbn">ISBN: {{.ISBN}}</div>
                    <div class="authors">
                        Authors: 
                        {{range $i, $email := .Authors}}
                            {{if $i}}, {{end}}
                            {{$author := index $.AuthorsMap $email}}
                            {{if $author}}
                                <span class="author-name">{{$author.FirstName}} {{$author.LastName}}</span>
                            {{else}}
                                {{$email}}
                            {{end}}
                        {{end}}
                    </div>
                    <div class="published">Published: {{.PublishedAt.Format "02.01.2006"}}</div>
                </div>
                {{end}}
            </div>
        </div>
        {{end}}
    </div>
</body>
</html>
`