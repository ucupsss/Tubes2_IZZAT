import { NavLink } from "react-router-dom";

const links = [
  { to: "/", label: "Home", end: true },
  { to: "/traversal", label: "DOM Traversal" },
  { to: "/about", label: "About Us" },
];

export default function Navbar() {
  return (
    <nav className="navbar">
      <span className="brand"></span>
      {links.map((link) => (
        <NavLink
          key={link.to}
          to={link.to}
          end={link.end}
          className={({ isActive }) => (isActive ? "nav-button active" : "nav-button")}
        >
          {link.label}
        </NavLink>
      ))}
    </nav>
  );
}