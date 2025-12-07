export type OrderStatus = 'pending' | 'confirmed' | 'cancelled';
export const orderStatus: Record<OrderStatus, string> = {
	pending: 'Pending',
	confirmed: 'Diterima',
	cancelled: 'Dicancel'
};

export type Customer = {
	id: number;
	name: string;
	phone: string;
	address: string;
	created_at: string;
};

export type Product = {
	id: string;
	user_id: string;
	name: string;
	price: number;
	stock: number;
	image_url: string;
	public_id: string;
	created_at: string;
};

export type OrderItem = {
	id: string;
	order_id: string;
	product_id: string;
	product: Product;
	quantity: number;
	total_price: number;
	created_at: string;
};

export type Order = {
	id: string;
	user_id: string;
	customer_id: string;
	customer: Customer;
	order_items: OrderItem[];
	total_price: number;
	status: OrderStatus;
	order_date: string;
	created_at: string;
};

export type OrdersResponse = {
	message: string;
	data: Order[];
	meta: PaginationMeta;
};
