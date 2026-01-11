import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { orderAPI, couponAPI } from '../services/api';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';
import { Coupon } from '../types';

const Checkout = () => {
  const { cart } = useCart();
  const { user, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [address, setAddress] = useState('');
  const [paymentType, setPaymentType] = useState('credit_card');
  const [couponCode, setCouponCode] = useState('');
  const [appliedCoupon, setAppliedCoupon] = useState<Coupon | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    if (!cart || cart.items.length === 0) {
      navigate('/cart');
      return;
    }
    if (user?.shipping_address) {
      setAddress(user.shipping_address);
    }
  }, [isAuthenticated, cart, user, navigate]);

  const handleApplyCoupon = async () => {
    try {
      const coupons = await couponAPI.getAllCoupons();
      const coupon = coupons.coupons?.find((c: Coupon) => c.code === couponCode && c.active);
      if (coupon) {
        setAppliedCoupon(coupon);
        setError('');
      } else {
        setError('Invalid or inactive coupon code');
        setAppliedCoupon(null);
      }
    } catch (error) {
      setError('Failed to apply coupon');
      setAppliedCoupon(null);
    }
  };

  const calculateDiscount = () => {
    if (!appliedCoupon || !cart) return 0;
    if (appliedCoupon.is_percent) {
      return (cart.total_price * appliedCoupon.discount) / 100;
    }
    return appliedCoupon.discount;
  };

  const calculateTotal = () => {
    if (!cart) return 0;
    return cart.total_price - calculateDiscount();
  };

  const handleCheckout = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!address.trim()) {
      setError('Please enter a shipping address');
      return;
    }
    if (!cart) {
      setError('Cart is empty');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const userId = localStorage.getItem('userId') || '';
      const response = await orderAPI.checkout({
        user_id: userId,
        cart_id: cart.cart_id,
        payment_type: paymentType,
        address: address,
        coupon_code: appliedCoupon?.code || '',
      });
      alert(`Order placed successfully! Expected delivery: ${response.expected_delivery}`);
      navigate('/orders');
    } catch (error: any) {
      setError(error.response?.data?.error || 'Checkout failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  if (!cart || cart.items.length === 0) {
    return null;
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-6">Checkout</h1>
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Checkout Form */}
        <div className="lg:col-span-2">
          <form onSubmit={handleCheckout} className="bg-white rounded-lg shadow-md p-6 space-y-6">
            {error && (
              <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
                {error}
              </div>
            )}

            {/* Shipping Address */}
            <div>
              <h2 className="text-xl font-bold mb-4">Shipping Address</h2>
              <textarea
                value={address}
                onChange={(e) => setAddress(e.target.value)}
                rows={4}
                required
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                placeholder="Enter your shipping address"
              />
            </div>

            {/* Payment Method */}
            <div>
              <h2 className="text-xl font-bold mb-4">Payment Method</h2>
              <select
                value={paymentType}
                onChange={(e) => setPaymentType(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
              >
                <option value="credit_card">Credit Card</option>
                <option value="debit_card">Debit Card</option>
                <option value="paypal">PayPal</option>
                <option value="cash_on_delivery">Cash on Delivery</option>
              </select>
            </div>

            {/* Coupon Code */}
            <div>
              <h2 className="text-xl font-bold mb-4">Coupon Code (Optional)</h2>
              <div className="flex space-x-2">
                <input
                  type="text"
                  value={couponCode}
                  onChange={(e) => setCouponCode(e.target.value)}
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                  placeholder="Enter coupon code"
                />
                <button
                  type="button"
                  onClick={handleApplyCoupon}
                  className="bg-gray-200 hover:bg-gray-300 px-6 py-2 rounded-md transition"
                >
                  Apply
                </button>
              </div>
              {appliedCoupon && (
                <p className="mt-2 text-green-600">
                  Coupon applied: {appliedCoupon.code} -{' '}
                  {appliedCoupon.is_percent
                    ? `${appliedCoupon.discount}% off`
                    : `$${appliedCoupon.discount} off`}
                </p>
              )}
            </div>

            {/* Order Items Summary */}
            <div>
              <h2 className="text-xl font-bold mb-4">Order Items</h2>
              <div className="space-y-2">
                {cart.items.map((item) => (
                  <div key={item.product_id} className="flex justify-between text-sm">
                    <span>
                      {item.product_name} x {item.quantity}
                    </span>
                    <span>${item.subtotal.toFixed(2)}</span>
                  </div>
                ))}
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-amazon-orange hover:bg-amazon-dark text-white font-semibold py-3 px-6 rounded transition disabled:opacity-50"
            >
              {loading ? 'Processing...' : 'Place Order'}
            </button>
          </form>
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
              {appliedCoupon && (
                <div className="flex justify-between text-green-600">
                  <span>Discount:</span>
                  <span>-${calculateDiscount().toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between text-sm text-gray-600">
                <span>Shipping:</span>
                <span>Free</span>
              </div>
            </div>
            <div className="border-t pt-4">
              <div className="flex justify-between text-lg font-bold">
                <span>Total:</span>
                <span className="text-amazon-orange">${calculateTotal().toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Checkout;
