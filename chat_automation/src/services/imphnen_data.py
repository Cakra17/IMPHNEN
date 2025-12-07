from dataclasses import dataclass
from typing import List

@dataclass
class AddCustomerData:
  id: int
  name: str
  address: str
  phone: str

@dataclass
class UpdateCustomerData:
  name: str
  address: str
  phone: str

@dataclass
class OrderItems:
  product_id: str
  quantity: int

@dataclass
class CreateOrdersData:
  customer_id: int
  items: List[OrderItems]
  merchant_id: str

@dataclass
class AddMerchantData:
  id: str
  name: str
  address: str
  phone: str

@dataclass
class UpdateMerchantData:
  name: str
  address: str
  phone: str

@dataclass
class OrderItemDetail:
  product_name: str
  quantity: int

@dataclass
class ListCustomerOrder:
  status: str
  total_price: int
  user_id: str
  order_date: str
  items: List[OrderItemDetail]