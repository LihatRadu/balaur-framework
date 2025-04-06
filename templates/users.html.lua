<!DOCTYPE html>
<html>
<head>
    <title><%= title %></title>
</head>
<body>
    <h1><%= title %></h1>
    <table>
        <thead>
            <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
            </tr>
        </thead>
        <tbody>
            <% for _, user in ipairs(users) do %>
                <tr>
                    <td><%= user.ID %></td>
                    <td><%= user.Username %></td>
                    <td><%= user.Email %></td>
                </tr>
            <% end %>
        </tbody>
    </table>
</body>
</html>
