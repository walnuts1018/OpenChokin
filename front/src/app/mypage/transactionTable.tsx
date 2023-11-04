import { Checkbox } from "@mui/material";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useLayoutEffect, useRef } from "react";
import { MoneyTransaction } from "./type";
import { ThemeProvider, createTheme, styled } from "@mui/material/styles";
import { useState, useEffect } from "react";

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  border: 0,
  "&:nth-of-type(odd)": {
    backgroundColor: theme.palette.action.hover,
  },
}));

const StyledTableCell = styled(TableCell)(({ theme }) => ({
  borderWidth: "0 1px",
}));

export function TransactionTable({
  moneyPoolID,
  scroll,
}: {
  moneyPoolID: string;
  scroll: boolean;
}) {
  const [transactions, setTransactions] = useState<MoneyTransaction[]>([]);
  const scrollBottomRef = useRef<HTMLTableRowElement>(null);

  useEffect(() => {
    const getMoneyPools = async () => {
      const res = await fetch("/backend/v1/moneypools?type=summary", {
        method: "GET",
      });
      if (res.ok) {
        const mps: MoneyPoolSum[] = await res.json();
        setMoneyPoolSums(mps);
      }
    };
    getMoneyPools();
  }, []);

  useLayoutEffect(() => {
    if (scroll) {
      {
        /*scrollBottomRef?.current?.scrollIntoView();*/
      }
    }
  }, [scroll]);
  return (
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
          {transactions.map((transaction) => (
            <StyledTableRow key={transaction.id}>
              <StyledTableCell
                align="center"
                component="th"
                scope="row"
                className="border-l-0"
              >
                {transaction.date.toLocaleDateString()}
              </StyledTableCell>
              <StyledTableCell align="left" className="">
                {transaction.title}
              </StyledTableCell>
              <StyledTableCell align="right" className=" border-r-0">
                {transaction.amount.toLocaleString(undefined, {
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
  );
}
