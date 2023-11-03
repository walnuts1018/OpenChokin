import { Checkbox } from "@mui/material";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { ThemeProvider, createTheme, styled } from "@mui/material/styles";
import { useLayoutEffect, useRef } from "react";

type MoneyTransaction = {
  id: number;
  date: Date;
  title: string;
  amount: number;
  version: number;
};

let rows: MoneyTransaction[] = [];

for (let i = 0; i < 100; i++) {
  rows.push({
    id: i,
    date: new Date(),
    title: "ごはん",
    amount: 1000,
    version: 1,
  });
}

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  border: 0,
  "&:nth-of-type(odd)": {
    backgroundColor: theme.palette.action.hover,
  },
}));

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  borderWidth: "0 1px",
}));

const theme = createTheme({
  palette: {
    primary: {
      main: "#f9842c",
      dark: "#FA6C28",
    },
  },
  typography: {
    fontFamily: "var(--font-Noto)",
  },
});

export function TransactionTable() {
  const scrollBottomRef = useRef<HTMLTableRowElement>(null);

  useLayoutEffect(() => {
    scrollBottomRef?.current?.scrollIntoView();
  }, []);
  return (
    <ThemeProvider theme={theme}>
      <TableContainer>
        <Table sx={{ minWidth: 200 }} size="small" stickyHeader>
          <TableHead>
            <TableRow sx={{ fontWeight: 700 }}>
              <StyledTableCell
                align="center"
                className="border-l-0 w-2/12 border-b-2"
              >
                日付
              </StyledTableCell>
              <StyledTableCell align="center" className="border-b-2">
                タイトル
              </StyledTableCell>
              <StyledTableCell
                align="center"
                className="border-r-0 w-40 border-b-2"
              >
                金額&nbsp;(円)
              </StyledTableCell>
              {/*
              <StyledTableCell align="center" className="border-r-0 w-20">
                公開
              </StyledTableCell>
              */}
            </TableRow>
          </TableHead>
          <TableBody className="h-full w-full">
            {rows.map((row) => (
              <StyledTableRow key={row.id}>
                <StyledTableCell
                  align="center"
                  component="th"
                  scope="row"
                  className="border-l-0"
                >
                  {row.date.toLocaleDateString()}
                </StyledTableCell>
                <StyledTableCell align="left" className="">
                  {row.title}
                </StyledTableCell>
                <StyledTableCell align="right" className=" border-r-0">
                  {row.amount.toLocaleString(undefined, {
                    maximumFractionDigits: 5,
                  })}
                </StyledTableCell>
                {/*
                <StyledTableCell align="center" className=" border-r-0 ">
                  <Checkbox size="small" className="p-0" />
                </StyledTableCell>
                */}
              </StyledTableRow>
            ))}
            <tr className="h-0" ref={scrollBottomRef} />
          </TableBody>
        </Table>
      </TableContainer>
    </ThemeProvider>
  );
}
