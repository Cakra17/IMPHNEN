import requests
from . import imphnen_data

class ImphnenService:
  base_url: str
  def __init__(self, base_url: str): 
    self.base_url = base_url

  def addCustomer(self, data: imphnen_data.AddCustomerData) -> str:
    try:
      headers = {
        "Content-Type": "application/json",
      }

      payload = {
        "id": data.id,
        "name": data.name,
        "address": data.address,
        "phone": data.phone
      }
      
      resp = requests.post(f"{self.base_url}/customers", json=payload, headers=headers, timeout=10)

      resp.raise_for_status()

      datas = resp.json()

      return datas["message"]
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal menambahkan customer.")
        except:
          return f"Gagal menambahkan customer. Status: {e.response.status_code}"
      return "Gagal menambahkan customer."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def getCustomerByID(self, id: int):
    try:
      headers = {
        "Content-Type": "application/json",
      }
      resp = requests.get(f"{self.base_url}/customers/{id}", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        if e.response.status_code == 404:
          return "Data customer tidak ditemukan."
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal mendapatkan data customer.")
        except:
          return f"Gagal mendapatkan data customer. Status: {e.response.status_code}"
      return "Gagal mendapatkan data customer."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."
    
  def updateCustomer(self, id: int, data: imphnen_data.UpdateCustomerData):
    try:
      headers = {
        "Content-Type": "application/json",
      }

      payload = {
        "name": data.name,
        "address": data.address,
        "phone": data.phone
      }
      resp = requests.put(f"{self.base_url}/customers/{id}", json=payload, headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas["message"]
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal memperbarui customer.")
        except:
          return f"Gagal memperbarui customer. Status: {e.response.status_code}"
      return "Gagal memperbarui customer."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def deleteCustomer(self, id: int):
    try:
      headers = {
        "Content-Type": "application/json",
      }
 
      resp = requests.delete(f"{self.base_url}/customers/{id}", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal menghapus customer.")
        except:
          return f"Gagal menghapus customer. Status: {e.response.status_code}"
      return "Gagal menghapus customer."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def listCustomerOrders(self, customerID: int):
    try:
      headers = {
        "Content-Type": "application/json",
      }
 
      resp = requests.get(f"{self.base_url}/customers/{customerID}/orders", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal mendapatkan daftar pesanan.")
        except:
          return f"Gagal mendapatkan daftar pesanan. Status: {e.response.status_code}"
      return "Gagal mendapatkan daftar pesanan."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def listProductByMerchant(self, merchantID):
    try:
      headers = {
        "Content-Type": "application/json",
      }
 
      resp = requests.get(f"{self.base_url}/merchants/{merchantID}/products", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas["data"]
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal mendapatkan daftar produk.")
        except:
          return f"Gagal mendapatkan daftar produk. Status: {e.response.status_code}"
      return "Gagal mendapatkan daftar produk."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def createOrder(self, data: imphnen_data.CreateOrdersData):
    try:
      headers = {
        "Content-Type": "application/json",
      }
      
      payload = {
        "customer_id": data.customer_id,
        "items": [{"product_id": item.product_id, "quantity": item.quantity} for item in data.items],
        "merchant_id": data.merchant_id
      }
 
      resp = requests.post(f"{self.base_url}/orders", json=payload, headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas["data"]
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal membuat pesanan.")
        except:
          return f"Gagal membuat pesanan. Status: {e.response.status_code}"
      return "Gagal membuat pesanan."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def deleteOrder(self, orderID: str):
    try:
      headers = {
        "Content-Type": "application/json",
      }
 
      resp = requests.delete(f"{self.base_url}/orders/{orderID}", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal menghapus pesanan.")
        except:
          return f"Gagal menghapus pesanan. Status: {e.response.status_code}"
      return "Gagal menghapus pesanan."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def cancelOrder(self, orderID):
    try:
      headers = {
        "Content-Type": "application/json",
      }
 
      resp = requests.patch(f"{self.base_url}/orders/{orderID}/cancel", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()
      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal membatalkan pesanan.")
        except:
          return f"Gagal membatalkan pesanan. Status: {e.response.status_code}"
      return "Gagal membatalkan pesanan."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi." 
  
  def getAllMerchant(self):
    try:
      headers = {
        "Content-Type": "application/json",
      }
      resp = requests.get(f"{self.base_url}/merchants", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas["data"]
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal mendapatkan daftar merchant.")
        except:
          return f"Gagal mendapatkan daftar merchant. Status: {e.response.status_code}"
      return "Gagal mendapatkan daftar merchant."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi."

  def getMerchantProductByID(self, id: str):
    try:
      headers = {
        "Content-Type": "application/json",
      }
      resp = requests.get(f"{self.base_url}/merchants/{id}/products", headers=headers, timeout=10)
      resp.raise_for_status()

      datas = resp.json()

      return datas
    except requests.exceptions.Timeout:
      return "Koneksi timeout. Silakan coba lagi."
    except requests.exceptions.ConnectionError:
      return "Gagal terhubung ke server. Silakan coba lagi."
    except requests.exceptions.HTTPError as e:
      if e.response is not None:
        try:
          error_data = e.response.json()
          return error_data.get("message", "Gagal mendapatkan produk merchant.")
        except:
          return f"Gagal mendapatkan produk merchant. Status: {e.response.status_code}"
      return "Gagal mendapatkan produk merchant."
    except Exception:
      return "Terjadi kesalahan. Silakan coba lagi." 