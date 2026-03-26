export default function Navigation() {
  const navItems = [
    { href: '/matches', label: 'Matches' },
    { href: '/me', label: 'My Profile' },
  ];

  return (
    <nav className="flex flex-col space-y-2">
      {navItems.map((item) => (
        <a
          key={item.href}
          href={item.href}
          className="flex items-center space-x-2 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"
        >
          <span>{item.label}</span>
        </a>
      ))}
    </nav>
  );
}
