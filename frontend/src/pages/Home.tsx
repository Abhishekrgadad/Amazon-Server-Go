import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { productAPI } from '../services/api';
import { Product } from '../types';
import { useCart } from '../context/CartContext';
import { useAuth } from '../context/AuthContext';

const Home = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [searchTerm, setSearchTerm] = useState('');
  const [filters, setFilters] = useState({
    category: '',
    brand: '',
    minPrice: '',
    maxPrice: '',
  });
  const { addToCart } = useCart();
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    fetchProducts();
  }, [page]);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      const response = await productAPI.getActiveProducts(page);
      setProducts(response.products || []);
    } catch (error) {
      console.error('Error fetching products:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    try {
      setLoading(true);
      const filterParams: any = {};
      if (searchTerm) filterParams.name = searchTerm;
      if (filters.category) filterParams.category = filters.category;
      if (filters.brand) filterParams.brand = filters.brand;
      if (filters.minPrice) filterParams.min_price = parseFloat(filters.minPrice);
      if (filters.maxPrice) filterParams.max_price = parseFloat(filters.maxPrice);

      const response = await productAPI.filterProducts(filterParams);
      setProducts(response.products || []);
    } catch (error) {
      console.error('Error searching products:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = async (productId: string) => {
    if (!isAuthenticated) {
      alert('Please login to add items to cart');
      return;
    }
    try {
      await addToCart(productId, 1);
      alert('Item added to cart!');
    } catch (error: any) {
      alert(error.message || 'Failed to add to cart');
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Hero Banner */}
      <div className="bg-gradient-to-r from-amazon-orange to-amazon-yellow rounded-lg p-8 mb-8 text-white">
        <h1 className="text-4xl font-bold mb-4">Welcome to Amazon Clone</h1>
        <p className="text-xl">Shop the best deals on thousands of products</p>
      </div>

      {/* Search and Filters */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="Search products..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
          />
          <input
            type="text"
            placeholder="Category"
            value={filters.category}
            onChange={(e) => setFilters({ ...filters, category: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
          />
          <input
            type="text"
            placeholder="Brand"
            value={filters.brand}
            onChange={(e) => setFilters({ ...filters, brand: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
          />
          <button
            onClick={handleSearch}
            className="bg-amazon-orange text-white px-6 py-2 rounded-md hover:bg-amazon-dark transition"
          >
            Search
          </button>
        </div>
        <div className="grid grid-cols-2 gap-4 mt-4">
          <input
            type="number"
            placeholder="Min Price"
            value={filters.minPrice}
            onChange={(e) => setFilters({ ...filters, minPrice: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
          />
          <input
            type="number"
            placeholder="Max Price"
            value={filters.maxPrice}
            onChange={(e) => setFilters({ ...filters, maxPrice: e.target.value })}
            className="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-amazon-orange"
          />
        </div>
      </div>

      {/* Products Grid */}
      {loading ? (
        <div className="text-center py-12">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-amazon-orange"></div>
          <p className="mt-4 text-gray-600">Loading products...</p>
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600 text-xl">No products found</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {products.map((product) => (
            <div
              key={product._id}
              className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition"
            >
              <Link to={`/product/${product._id}`}>
                <div className="aspect-w-1 aspect-h-1 bg-gray-200 p-4">
                  <div className="w-full h-48 bg-gray-100 flex items-center justify-center">
                    <span className="text-gray-400">No Image</span>
                  </div>
                </div>
              </Link>
              <div className="p-4">
                <Link to={`/product/${product._id}`}>
                  <h3 className="font-semibold text-lg mb-2 line-clamp-2 hover:text-amazon-orange transition">
                    {product.name}
                  </h3>
                </Link>
                <div className="flex items-center mb-2">
                  {product.average_rating ? (
                    <>
                      <span className="text-yellow-400">★</span>
                      <span className="ml-1 text-sm">{product.average_rating.toFixed(1)}</span>
                      <span className="ml-2 text-sm text-gray-500">
                        ({product.total_reviews || 0} reviews)
                      </span>
                    </>
                  ) : (
                    <span className="text-sm text-gray-500">No ratings yet</span>
                  )}
                </div>
                <p className="text-2xl font-bold text-amazon-orange mb-2">
                  ${product.price.toFixed(2)}
                </p>
                <p className="text-sm text-gray-600 mb-4 line-clamp-2">{product.description}</p>
                <button
                  onClick={() => handleAddToCart(product._id)}
                  className="w-full bg-amazon-yellow hover:bg-amazon-orange text-amazon-dark font-semibold py-2 px-4 rounded transition"
                >
                  Add to Cart
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {!loading && products.length > 0 && (
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

export default Home;
