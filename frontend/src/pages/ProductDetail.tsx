import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { productAPI, reviewAPI } from '../services/api';
import { Product, Review } from '../types';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

const ProductDetail = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [quantity, setQuantity] = useState(1);
  const [reviewForm, setReviewForm] = useState({
    rating: 5,
    comment: '',
    order_id: '',
  });
  const { addToCart } = useCart();
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    if (id) {
      fetchProduct();
      fetchReviews();
    }
  }, [id]);

  const fetchProduct = async () => {
    try {
      setLoading(true);
      const data = await productAPI.getProductById(id!);
      setProduct(data);
    } catch (error) {
      console.error('Error fetching product:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchReviews = async () => {
    try {
      const data = await reviewAPI.getProductReviews(id!);
      setReviews(data.reviews || []);
    } catch (error) {
      console.error('Error fetching reviews:', error);
    }
  };

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      alert('Please login to add items to cart');
      navigate('/login');
      return;
    }
    try {
      await addToCart(id!, quantity);
      alert('Item added to cart!');
    } catch (error: any) {
      alert(error.message || 'Failed to add to cart');
    }
  };

  const handleSubmitReview = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isAuthenticated) {
      alert('Please login to submit a review');
      return;
    }
    try {
      await reviewAPI.addReview({
        product_id: id!,
        order_id: reviewForm.order_id,
        rating: reviewForm.rating,
        comment: reviewForm.comment,
      });
      alert('Review submitted successfully!');
      setReviewForm({ rating: 5, comment: '', order_id: '' });
      fetchReviews();
    } catch (error: any) {
      alert(error.response?.data?.error || 'Failed to submit review');
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-amazon-orange"></div>
          <p className="mt-4 text-gray-600">Loading product...</p>
        </div>
      </div>
    );
  }

  if (!product) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <p className="text-gray-600 text-xl">Product not found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div className="bg-white rounded-lg shadow-lg overflow-hidden">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 p-8">
          {/* Product Image */}
          <div className="bg-gray-100 rounded-lg p-8 flex items-center justify-center">
            <div className="w-full h-96 bg-gray-200 flex items-center justify-center">
              <span className="text-gray-400 text-xl">Product Image</span>
            </div>
          </div>

          {/* Product Info */}
          <div>
            <h1 className="text-3xl font-bold mb-4">{product.name}</h1>
            <div className="flex items-center mb-4">
              {product.average_rating ? (
                <>
                  <div className="flex text-yellow-400">
                    {'★'.repeat(Math.floor(product.average_rating))}
                  </div>
                  <span className="ml-2 text-lg font-semibold">{product.average_rating.toFixed(1)}</span>
                  <span className="ml-2 text-gray-600">
                    ({product.total_reviews || 0} reviews)
                  </span>
                </>
              ) : (
                <span className="text-gray-600">No ratings yet</span>
              )}
            </div>
            <div className="border-t border-b py-4 my-4">
              <p className="text-4xl font-bold text-amazon-orange mb-2">
                ${product.price.toFixed(2)}
              </p>
              <p className="text-sm text-gray-600">
                In Stock: {product.stock_quantity} available
              </p>
            </div>
            <p className="text-gray-700 mb-6">{product.description}</p>
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Quantity:
              </label>
              <select
                value={quantity}
                onChange={(e) => setQuantity(parseInt(e.target.value))}
                className="border border-gray-300 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-amazon-orange"
              >
                {[...Array(Math.min(product.stock_quantity, 10))].map((_, i) => (
                  <option key={i + 1} value={i + 1}>
                    {i + 1}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-4">
              <button
                onClick={handleAddToCart}
                className="w-full bg-amazon-yellow hover:bg-amazon-orange text-amazon-dark font-semibold py-3 px-6 rounded transition"
              >
                Add to Cart
              </button>
              <button
                onClick={() => {
                  handleAddToCart();
                  navigate('/checkout');
                }}
                className="w-full bg-amazon-orange hover:bg-amazon-dark text-white font-semibold py-3 px-6 rounded transition"
              >
                Buy Now
              </button>
            </div>
            <div className="mt-6 text-sm text-gray-600">
              <p><strong>Brand:</strong> {product.brand}</p>
              <p><strong>Category:</strong> {product.category}</p>
            </div>
          </div>
        </div>

        {/* Reviews Section */}
        <div className="border-t p-8">
          <h2 className="text-2xl font-bold mb-6">Customer Reviews</h2>
          
          {/* Review Form */}
          {isAuthenticated && (
            <form onSubmit={handleSubmitReview} className="mb-8 bg-gray-50 p-6 rounded-lg">
              <h3 className="text-lg font-semibold mb-4">Write a Review</h3>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Order ID (required for review)
                  </label>
                  <input
                    type="text"
                    value={reviewForm.order_id}
                    onChange={(e) => setReviewForm({ ...reviewForm, order_id: e.target.value })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                    placeholder="Enter your order ID"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Rating
                  </label>
                  <select
                    value={reviewForm.rating}
                    onChange={(e) => setReviewForm({ ...reviewForm, rating: parseInt(e.target.value) })}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                  >
                    {[5, 4, 3, 2, 1].map((rating) => (
                      <option key={rating} value={rating}>
                        {rating} {rating === 1 ? 'Star' : 'Stars'}
                      </option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Comment
                  </label>
                  <textarea
                    value={reviewForm.comment}
                    onChange={(e) => setReviewForm({ ...reviewForm, comment: e.target.value })}
                    rows={4}
                    className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
                    placeholder="Share your experience with this product..."
                    required
                  />
                </div>
                <button
                  type="submit"
                  className="bg-amazon-orange text-white px-6 py-2 rounded-md hover:bg-amazon-dark transition"
                >
                  Submit Review
                </button>
              </div>
            </form>
          )}

          {/* Reviews List */}
          <div className="space-y-6">
            {reviews.length === 0 ? (
              <p className="text-gray-600">No reviews yet. Be the first to review!</p>
            ) : (
              reviews.map((review) => (
                <div key={review._id} className="border-b pb-4">
                  <div className="flex items-center mb-2">
                    <div className="flex text-yellow-400">
                      {'★'.repeat(review.rating)}
                    </div>
                    <span className="ml-2 text-sm text-gray-600">
                      {review.created_at ? new Date(review.created_at).toLocaleDateString() : ''}
                    </span>
                  </div>
                  <p className="text-gray-700">{review.comment}</p>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ProductDetail;
