<!DOCTYPE html>
<html>
<head>
    <title><%= title %></title>
</head>
<body>
    <h1><%= greeting %></h1>
    <ul>
        <% for i, item in ipairs(items) do %>
            <li><%= item %></li>
        <% end %>
    </ul>
</body>
</html>
