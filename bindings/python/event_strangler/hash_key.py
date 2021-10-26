from typing import Optional

import bindings


class HashKeyOptions:
    def __init__(self, name: Optional[str], expression: Optional[str] = None):
        self.name: Optional[str] = name
        self.expression: Optional[str] = expression

    @staticmethod
    def to_binding(model) -> bindings.HashKeyOptions:
        result = bindings.HashKeyOptions()
        result.Name = model.name
        result.Expression = model.name
        return result
