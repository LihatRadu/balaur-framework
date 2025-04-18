local Button = {}

function Button:new(text, onClick)
	local btn = {
		text = text,
		onClick = onClick,
	}
	setmetatable(btn, { __index = Button })
	return btn
end

function Button:render()
	return string.format(
		[[
      <button onclick=%s>%s</button>
    ]],
		self.onClick,
		self.text
	)
end

return Button
