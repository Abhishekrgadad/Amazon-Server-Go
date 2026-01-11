import { useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

const Cart = () => {
  const { cart, loading, updateCartItem, removeFromCart, refreshCart } = useCart();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      refreshCart();
    }
  }, [isAuthenticated]);

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600 text-xl mb-4">Please login to view your cart</p>
          <Link
            to="/login"
            className="bg-amazon-orange text-white px-6 py-2 rounded-md hover:bg-amazon-dark transition"
          >
            Login
          </Link>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-amazon-orange"></div>
          <p className="mt-4 text-gray-600">Loading cart...</p>
        </div>
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <h2 className="text-2xl font-bold mb-4">Your Amazon Cart is empty</h2>
          <p className="text-gray-600 mb-6">Shop today's deals</p>
          <Link
            to="/"
            className="bg-amazon-orange text-white px-6 py-2 rounded-md hover:bg-amazon-dark transition inline-block"
          >
            Continue Shopping
          </Link>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-6">Shopping Cart</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Cart Items */}
        <div className="lg:col-span-2 space-y-4">
          {cart.items.map((item) => (
            <div
              key={item.product_id}
              className="bg-white rounded-lg shadow-md p-6 flex items-start space-x-4"
            >
              <div className="w-24 h-24 bg-gray-200 rounded flex items-center justify-center flex-shrink-0">
                <span className="text-gray-400 text-xs">Image</span>
              </div>
              <div className="flex-1">
                <Link
                  to={`/product/${item.product_id}`}
                  className="text-lg font-semibold hover:text-amazon-orange transition"
                >
                  {item.product_name}
                </Link>
                <p className="text-amazon-orange font-bold text-xl mt-2">
                  ${item.price.toFixed(2)}
                </p>
                <div className="flex items-center space-x-4 mt-4">
                  <label className="text-sm font-medium">Quantity:</label>
                  <select
                    value={item.quantity}
                    onChange={(e) => updateCartItem(item.product_id, parseInt(e.target.value))}
                    className="border border-gray-300 rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                  >
                    {[...Array(10)].map((_, i) => (
                      <option key={i + 1} value={i + 1}>
                        {i + 1}
                      </option>
                    ))}
                  </select>
                  <button
                    onClick={() => removeFromCart(item.product_id)}
                    className="text-red-600 hover:text-red-800 text-sm"
                  >
                    Delete
                  </button>
                </div>
                <p className="text-gray-600 mt-2">
                  Subtotal: <span className="font-semibold">${item.subtotal.toFixed(2)}</span>
                </p>
              </div>
            </div>
          ))}
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-lg shadow-md p-6 sticky top-4">
            <h2 className="text-xl font-bold mb-4">Order Summary</h2>
            <div className="space-y-2 mb-4">
              <div className="flex justify-between">
                <span>Subtotal ({cart.items.length} items):</span>
                <span className="font-semibold">${cart.total_price.toFixed(2)}</span>
              </div>
              <div className="flex justify-between text-sm text-gray-600">
                <span>Shipping:</span>
                <span>Free</span>
              </div>
            </div>
            <div className="border-t pt-4 mb-4">
              <div className="flex justify-between text-lg font-bold">
                <span>Total:</span>
                <span className="text-amazon-orange">${cart.total_price.toFixed(2)}</span>
              </div>
            </div>
            <button
              onClick={() => navigate('/checkout')}
              className="w-full bg-amazon-yellow hover:bg-amazon-orange text-amazon-dark font-semibold py-3 px-6 rounded transition"
            >
              Proceed to Checkout
            </button>
            <Link
              to="/"
              className="block text-center text-amazon-orange hover:text-amazon-dark mt-4"
            >
              Continue Shopping
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Cart;
