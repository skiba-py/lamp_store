import { Box, Image, Text, IconButton, Flex } from '@chakra-ui/react';
import { FaShoppingBasket } from 'react-icons/fa';
import { Link } from 'react-router-dom';

export default function ProductCard({ product, onAddToCart }) {
  return (
    <Box
      w="100%"
      bg="white"
      borderRadius="2xl"
      boxShadow="md"
      p={4}
      display="flex"
      flexDirection="column"
      transition="box-shadow 0.2s"
      _hover={{ boxShadow: 'xl' }}
      minH="260px"
    >
      <Link to={`/product/${product.id}`} style={{ textDecoration: 'none' }}>
        <Box w="100%" aspectRatio={4/3} mb={2} overflow="hidden" borderRadius="lg">
          <Image
            src={product.image}
            alt={product.name}
            objectFit="cover"
            w="100%"
            h="100%"
          />
        </Box>
        <Text fontWeight="bold" fontSize="md" mb={1}>{product.name}</Text>
      </Link>
      <Flex align="center" justify="space-between">
        <Text fontWeight="bold" fontSize="lg">{product.price} ₽</Text>
        <IconButton
          icon={<FaShoppingBasket />}
          colorScheme="red"
          aria-label="В корзину"
          onClick={onAddToCart}
          size="sm"
          borderRadius="full"
        />
      </Flex>
    </Box>
  );
}
