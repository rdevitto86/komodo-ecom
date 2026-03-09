export interface CartItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
}

class CartState {
  #data = $state({
    items: [] as CartItem[],
    total: 0,
  })
  
  get items() { return this.#data.items; }
  get total() { return this.#data.total; }
}

export const cartState = new CartState();
