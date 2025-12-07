type Product = {
	created_at: string;
	id: string;
	image_url: string;
	name: string;
	price: number;
	public_id: string;
	stock: number;
	user_id: string;
};

type PaginationMeta = {
	page: number;
	per_page: number;
	total_data: number;
	total_page: number;
};
