import { Box, Flex, Text, Button, Spacer } from '@chakra-ui/react';
import { FaRegLightbulb } from 'react-icons/fa';
import { Link } from 'react-router-dom';

export default function Header() {
  return (
    <Box as="header" bg="white" py={4}>
      <Flex align="center" maxW="1200px" mx="auto" px={4}>
        <Link to="/" style={{ display: 'flex', alignItems: 'center', textDecoration: 'none' }}>
          <Text fontWeight="bold" fontSize="xl">LAMPS</Text>
          <Box ml={2}><FaRegLightbulb /></Box>
        </Link>
        <Spacer />
        <Flex gap={4} align="center">
          <Link to="/cart">
            <Button variant="outline" size="sm">Корзина</Button>
          </Link>
        </Flex>
      </Flex>
    </Box>
  );
}
