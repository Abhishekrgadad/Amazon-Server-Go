import { useState, useEffect } from 'react';
import { orderAPI } from '../services/api';
import { Order } from '../types';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const Orders = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    fetchOrders();
  }, [isAuthenticated, page, navigate]);

  const fetchOrders = async () => {
    try {
      setLoading(true);
      const userId = localStorage.getItem('userId') || '';
      const response = await orderAPI.viewOrders(userId, page);
      setOrders(response.orders || []);
    } catch (error) {
      console.error('Error fetching orders:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCancelOrder = async (orderId: string) => {
    if (!window.confirm('Are you sure you want to cancel this order?')) {
      return;
    }
    try {
      const userId = localStorage.getItem('userId') || '';
      await orderAPI.cancelOrder(orderId, userId);
      alert('Order cancelled successfully. Refund will be initiated soon.');
      fetchOrders();
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to cancel order');
    }
  };

  const handleReturnOrder = async (orderId: string) => {
    if (!window.confirm('Are you sure you want to return this order?')) {
      return;
    }
    try {
      await orderAPI.returnOrder(orderId);
      alert('Return initiated successfully. Refund will be processed after we receive the product.');
      fetchOrders();
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to return order');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'confirmed':
        return 'bg-blue-100 text-blue-800';
      case 'shipped':
        return 'bg-purple-100 text-purple-800';
      case 'delivered':
        return 'bg-green-100 text-green-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-amazon-orange"></div>
          <p className="mt-4 text-gray-600">Loading orders...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="text-3xl font-bold mb-6">Your Orders</h1>
      {orders.length === 0 ? (
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <p className="text-gray-600 text-xl mb-4">You have no orders yet</p>
          <a
            href="/"
            className="bg-amazon-orange text-white px-6 py-2 rounded-md hover:bg-amazon-dark transition inline-block"
          >
            Start Shopping
          </a>
        </div>
      ) : (
        <div className="space-y-6">
          {orders.map((order) => (
            <div key={order._id} className="bg-white rounded-lg shadow-md p-6">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <p className="text-sm text-gray-600">Order Placed</p>
                  <p className="font-semibold">
                    {new Date(order.created_at).toLocaleDateString()}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">Order ID</p>
                  <p className="font-semibold text-sm">{order._id}</p>
                </div>
                <div>
                  <span
                    className={`px-3 py-1 rounded-full text-sm font-semibold ${getStatusColor(
                      order.status
                    )}`}
                  >
                    {order.status}
                  </span>
                </div>
              </div>

              <div className="border-t pt-4 space-y-4">
                {order.items.map((item, index) => (
                  <div key={index} className="flex items-start space-x-4">
                    <div className="w-20 h-20 bg-gray-200 rounded flex items-center justify-center flex-shrink-0">
                      <span className="text-gray-400 text-xs">Image</span>
                    </div>
                    <div className="flex-1">
                      <h3 className="font-semibold">{item.product_name}</h3>
                      <p className="text-sm text-gray-600">
                        Quantity: {item.quantity} × ${item.price.toFixed(2)}
                      </p>
                    </div>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4 mt-4 flex justify-between items-center">
                <div>
                  <p className="text-sm text-gray-600">Total Amount</p>
                  <p className="text-xl font-bold text-amazon-orange">
                    ${order.total_price.toFixed(2)}
                  </p>
                  {order.discount > 0 && (
                    <p className="text-sm text-green-600">
                      Discount: -${order.discount.toFixed(2)}
                    </p>
                  )}
                </div>
                <div className="space-x-2">
                  {order.status.toLowerCase() !== 'cancelled' &&
                    order.status.toLowerCase() !== 'delivered' && (
                      <button
                        onClick={() => handleCancelOrder(order._id)}
                        className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md transition"
                      >
                        Cancel Order
                      </button>
                    )}
                  {order.status.toLowerCase() === 'delivered' && (
                    <button
                      onClick={() => handleReturnOrder(order._id)}
                      className="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-md transition"
                    >
                      Return Order
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {orders.length > 0 && (
        <div className="flex justify-center mt-8 space-x-2">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="px-4 py-2 bg-white border border-gray-300 rounded-md disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            Previous
          </button>
          <span className="px-4 py-2">Page {page}</span>
          <button
            onClick={() => setPage(page + 1)}
            className="px-4 py-2 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
};

export default Orders;
