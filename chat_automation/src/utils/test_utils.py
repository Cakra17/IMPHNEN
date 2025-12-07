import pytest

from utils.utils import is_complex_json


@pytest.mark.parametrize("json_input, expected", [
    # 1. Valid Complex JSON (Expected: True)
    ('{"name": "Bob"}', True),           # Simple Dict
    ('{}', True),                        # Empty Dict
    ('[1, 2, 3]', True),                 # Simple List
    ('[]', True),                        # Empty List
    ('{"a": [1, 2]}', True),             # Nested structure
    
    # 2. Valid Primitive JSON (Expected: False - your function excludes these)
    ('123', False),                      # Integer
    ('12.5', False),                     # Float
    ('"Just a string"', False),          # Quoted String
    ('true', False),                     # Boolean
    ('null', False),                     # Null
    
    # 3. Invalid JSON (Expected: False)
    ('Hello World', False),              # Raw text
    ("{'a': 1}", False),                 # Single quotes (Invalid JSON syntax)
    ('{a: 1}', False),                   # Missing key quotes (Invalid JSON syntax)
    ('', False),                         # Empty string (Fails json.loads)
    ('{"a": 1,}', False),                # Trailing comma (Invalid JSON)
    (None, False),                       # None type input
    (12345, False),                      # Raw integer input (TypeError)
])
def test_is_complex_json(json_input, expected):
    assert is_complex_json(json_input) == expected  