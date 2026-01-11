import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useCart } from '../context/CartContext';

const Navbar = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const { getCartItemCount } = useCart();
  const navigate = useNavigate();
  const cartCount = getCartItemCount();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <nav className="bg-amazon-dark text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center">
            <span className="text-2xl font-bold text-amazon-yellow">amazon</span>
          </Link>

          {/* Search Bar */}
          <div className="flex-1 max-w-2xl mx-4 hidden md:flex">
            <div className="w-full flex">
              <input
                type="text"
                placeholder="Search products..."
                className="flex-1 px-4 py-2 text-black rounded-l-md focus:outline-none"
              />
              <button className="bg-amazon-yellow text-amazon-dark px-6 py-2 rounded-r-md hover:bg-amazon-orange transition">
                Search
              </button>
            </div>
          </div>

          {/* Right Side */}
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <Link to="/profile" className="text-sm hover:text-amazon-yellow transition">
                  <div>Hello, {user?.email?.split('@')[0] || 'User'}</div>
                  <div className="font-semibold">Account & Lists</div>
                </Link>
                <Link to="/orders" className="text-sm hover:text-amazon-yellow transition">
                  <div>Returns</div>
                  <div className="font-semibold">& Orders</div>
                </Link>
                <Link to="/cart" className="flex items-center space-x-1 hover:text-amazon-yellow transition">
                  <div className="relative">
                    <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
                    </svg>
                    {cartCount > 0 && (
                      <span className="absolute -top-1 -right-1 bg-amazon-orange text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
                        {cartCount}
                      </span>
                    )}
                  </div>
                  <span className="font-semibold">Cart</span>
                </Link>
                <button
                  onClick={handleLogout}
                  className="text-sm hover:text-amazon-yellow transition"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="text-sm hover:text-amazon-yellow transition">
                  <div>Hello, Sign in</div>
                  <div className="font-semibold">Account & Lists</div>
                </Link>
                <Link to="/register" className="text-sm hover:text-amazon-yellow transition px-4 py-2 border border-gray-400 rounded">
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>

        {/* Secondary Nav */}
        <div className="flex items-center space-x-6 py-2 text-sm">
          <Link to="/" className="hover:text-amazon-yellow transition">All</Link>
          <Link to="/" className="hover:text-amazon-yellow transition">Today's Deals</Link>
          <Link to="/" className="hover:text-amazon-yellow transition">Customer Service</Link>
          <Link to="/" className="hover:text-amazon-yellow transition">Registry</Link>
          <Link to="/" className="hover:text-amazon-yellow transition">Gift Cards</Link>
          <Link to="/" className="hover:text-amazon-yellow transition">Sell</Link>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
