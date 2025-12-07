import json

def is_complex_json(json_str):
    """
    Checks if a string is valid JSON and if the root element is a dictionary (object) 
    or a list (array). Primitives (numbers, strings, booleans) return False.
    
    Args:
        json_str: The string to validate. Can be None or a non-string type.
        
    Returns:
        True if the string is valid JSON and its root is a dict or list, otherwise False.
    """
    if json_str is None:
        return False

    try:
        data = json.loads(json_str)
        # Only return True if the result is a dict or a list
        if isinstance(data, (dict, list)):
            return True
        
        return False
    except (ValueError, TypeError):
        return False