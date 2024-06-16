"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import { Container, Navbar, Nav } from "react-bootstrap";
function SmartLink({ children, path, part = false }) {
  const pathname = usePathname();
  let c = "nav-link ";
  if (pathname === path || (part && pathname.indexOf(path) != -1)) {
    c += "active ";
  }
  return (
    <li className="nav-item">
      <Link className={c} href={path}>
        {children}
      </Link>
    </li>
  );
}
export function Links() {
  const pathname = usePathname();

  return (
    <Navbar
      collapseOnSelect
      expand="md"
      className={pathname == "/" ? "bg-dark" : "bg-body-tertiary"}
    >
      <Container>
        <Navbar.Brand>
          <Image
            src="/gool.png"
            alt="Logo"
            width={36}
            height={30}
            className="d-inline-block align-text-top"
          />
        </Navbar.Brand>
        <Navbar.Brand>КОПИРАЙТРШТНАЯ</Navbar.Brand>

        <Navbar.Toggle aria-controls="responsive-navbar-nav" />

        <Navbar.Collapse id="responsive-navbar-nav">
          <Nav className="me-auto">
            <SmartLink path="/">ООООооо</SmartLink>
            <SmartLink part path="/pages">
              Галерея
            </SmartLink>
            <SmartLink path="/load">Загрузить</SmartLink>
            <SmartLink path="/about">Команда</SmartLink>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}
