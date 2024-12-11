-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS `bookstore`;

-- Create reviews table
CREATE TABLE IF NOT EXISTS `reviews` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `isbn` varchar(50) NOT NULL,
  `title` varchar(255) NOT NULL,
  `author` varchar(255) NOT NULL,
  `genre` varchar(50) NOT NULL,
  `descriptions` text NOT NULL,
  `price` decimal(5,2) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `uid` int(11) NOT NULL,
  `rating` int(11) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT 0
);

-- Create users table 

CREATE TABLE IF NOT EXISTS `users` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `email` varchar(100) UNIQUE NOT NULL,
  `password` varchar(100) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `login_token` varchar(100) DEFAULT NULL,
  `active` tinyint(1) DEFAULT 0,
  `activation_token` varchar(100) DEFAULT NULL
);

--  Create forget_passw table 

CREATE TABLE IF NOT EXISTS `forget_passw` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `uid` int(11) NOT NULL,
  `uri` varchar(100) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);

--  Create user_log table 

CREATE TABLE IF NOT EXISTS `user_log` (
  `id` int(11) PRIMARY KEY AUTO_INCREMENT NOT NULL,
  `activity` varchar(50) NOT NULL,
  `uid` int(11) NOT NULL,
  `created_at` datetime DEFAULT current_timestamp(),
  `superseded` tinyint(1) DEFAULT 0
);

--  Create books table 

CREATE TABLE IF NOT EXISTS `books` (
  `isbn` varchar(100) PRIMARY KEY NOT NULL,
  `title` varchar(255) NOT NULL,
  `author` varchar(255) NOT NULL,
  `genre` varchar(50) NOT NULL,
  `descriptions` text NOT NULL,
  `price` decimal(10,2) NOT NULL
);

-- Insert dummy data into users table
-- password user type will be password1 for both reset it after ward 

INSERT INTO `users` (`email`, `password`, `created_at`, `login_token`, `active`, `activation_token`) VALUES
('user1@example.com', '$2a$12$vChxDZJ8me0zA2gWMwq7MOcNYScff4xe6mIv/xEJNwfRDpVSXcure', current_timestamp(), NULL, 1, NULL),
('user2@example.com', '$2a$12$vChxDZJ8me0zA2gWMwq7MOcNYScff4xe6mIv/xEJNwfRDpVSXcure', current_timestamp(), NULL, 1, NULL);

-- Insert dummy data into books table 
INSERT INTO `books` (`isbn`, `title`, `author`, `genre`, `descriptions`, `price`) VALUES
('978-3-16-148410-0', 'Sapiens', 'Yoah N Harari', 'Reality', 'Human Kind Development', 19.99),
('978-1-23-456789-7', 'Animal Farm', 'George Orwell', 'Fiction', 'Politics & leadership', 29.99);

-- Insert dummy data into reviews table 
INSERT INTO `reviews` (`isbn`, `title`, `author`, `genre`, `descriptions`, `price`, `created_at`, `uid`, `rating`, `is_deleted`) VALUES
('978-3-16-148410-0', 'Book Title 1', 'Author 1', 'Fiction', 'Review for book 1', 19.99, current_timestamp(), 1, 5, 0),
('978-1-23-456789-7', 'Book Title 2', 'Author 2', 'Non-Fiction', 'Review for book 2', 29.99, current_timestamp(), 2, 4, 0);


--  Insert dummy data into user_log table 
INSERT INTO `user_log` (`activity`, `uid`, `created_at`, `superseded`) VALUES
('Login', 1, current_timestamp(), 0),
('Logout', 2, current_timestamp(), 0);
